package climber

import (
	"AutoArticle/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	arkutils "github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"golang.org/x/time/rate"
)

type Code int

const (
	SuccessCode Code = iota
	ServerOverloadedCode
	InternalServerError
	UnknownCode
	Accept
	numOfCodes
)
const (
	ServerName = "dream_qwen"
	PostPath   = "/api/v3/chat/completions"
)

func convertToCode(statusCode int) Code {
	switch {
	case statusCode >= http.StatusInternalServerError:
		return InternalServerError
	case statusCode == http.StatusTooManyRequests:
		return ServerOverloadedCode
	case statusCode == http.StatusOK:
		return SuccessCode
	default:
		return UnknownCode
	}
}

type Counter interface {
	Record(code Code)
	Count(code Code) int64
	Reset()
}

type counter struct {
	codes [numOfCodes]atomic.Int64
}

func NewCounter() Counter {
	return &counter{}
}

func (c *counter) Record(code Code) {
	_ = &c.codes[code]
	c.codes[code].Add(1)
}

func (c *counter) Count(code Code) int64 {
	_ = &c.codes[code]
	return c.codes[code].Load()
}

func (c *counter) Reset() {
	for code := range c.codes {
		c.codes[code].Store(0)
	}
}

type Option func(*Options)

func WithAdjustInterval(interval time.Duration) Option {
	return func(option *Options) {
		if interval >= time.Second {
			option.AdjustInterval = interval
		}
	}
}

func WithInitialRate(initial int64) Option {
	return func(option *Options) {
		if initial > 0 {
			option.InitialRate = initial
		}
	}
}

func WithMaxRate(max int64) Option {
	return func(option *Options) {
		if max > 0 {
			option.MaxRate = max
		}
	}
}

func defaultOptions() *Options {
	return &Options{
		AdjustInterval:   10 * time.Second, //跟新QPS的间隔
		InitialRate:      10,               // 初始qps速率
		MaxRate:          1000,             // 最大qps速率
		StrictErrorRatio: 0.05,             //错误比例  错误率阈值1，错误率高于该阈值，会降低请求频率。 对稳定性要求高时可降低该值。
		RelaxErrorRatio:  0.01,             // 允许少量错误时可适当提高该值。
		ProbeRatio:       0.05,             // 速率调整步长，波动大时可降低步长。
		FastProbeRatio:   0.1,

		RetryBaseDelay: 100 * time.Millisecond, // 单位 ms，初始重试延迟，根据服务响应速度调整。
		RetryMaxDelay:  3 * time.Second,        // 最大重试等待时间。
	}
}

type Options struct {
	AdjustInterval time.Duration

	InitialRate int64
	MaxRate     int64

	StrictErrorRatio float64
	RelaxErrorRatio  float64

	ProbeRatio     float64
	FastProbeRatio float64

	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

type QPSClimber struct {
	client    *arkruntime.Client
	limiter   *rate.Limiter
	counter   Counter
	options   *Options
	apiKey    string
	uniqueKey string
}

func NewQPSClimber(client *arkruntime.Client, apiKey string, opts ...Option) *QPSClimber {
	config := defaultOptions()
	for _, opt := range opts {
		opt(config)
	}
	uniqueKey := utils.GenerateUniqueId()
	climber := &QPSClimber{
		client:    client,
		limiter:   rate.NewLimiter(rate.Limit(config.InitialRate), int(config.InitialRate)),
		counter:   NewCounter(),
		options:   config,
		apiKey:    apiKey,
		uniqueKey: uniqueKey,
	}

	go climber.startAdjusting(context.Background())

	return climber
}

func (q *QPSClimber) DoRequest(ctx context.Context, req *model.CreateChatCompletionRequest) (model.ChatCompletionResponse, error) {
	backOff := newExponentialBackoff(q.options.RetryBaseDelay, q.options.RetryMaxDelay)

	for {
		select {
		case <-ctx.Done():
			return model.ChatCompletionResponse{}, ctx.Err()
		default:
		}

		// 检查是否有足够的令牌，未拿到令牌重试
		if !q.limiter.Allow() {
			backOff.wait()
			continue
		}

		q.Record(Accept)
		reqByte, _ := json.Marshal(req)
		println(string(reqByte))
		resp, err := q.client.CreateChatCompletion(ctx, req)

		//if resp.Usage.TotalTokens > 0 {
		//	util.SetTotalTokens(ctx, resp.Usage.TotalTokens, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, q.apiKey, resp.Model, "dbao")
		//
		//}
		if err == nil {
			return resp, nil
		}

		fmt.Printf("chat error: %v\n", err)
		statusCode, retry := needRetryError(err)
		q.Record(convertToCode(statusCode))

		if !retry {
			return resp, err
		}

		backOff.wait()
	}
}

func (q *QPSClimber) DoRequestStream(ctx context.Context, req *model.CreateChatCompletionRequest) (*arkutils.ChatCompletionStreamReader, error) {
	backOff := newExponentialBackoff(q.options.RetryBaseDelay, q.options.RetryMaxDelay)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 检查是否有足够的令牌，未拿到令牌重试
		if !q.limiter.Allow() {
			backOff.wait()
			continue
		}

		q.Record(Accept)
		stream, err := q.client.CreateChatCompletionStream(ctx, req)
		if err == nil {
			return stream, nil
		}

		fmt.Printf("stream chat error: %v\n", err)
		statusCode, retry := needRetryError(err)
		q.Record(convertToCode(statusCode))

		if !retry {
			return stream, err
		}

		backOff.wait()
	}
}

func (q *QPSClimber) Record(code Code) {
	q.counter.Record(code)
}

func (q *QPSClimber) Wait(ctx context.Context) error {
	return q.limiter.Wait(ctx)
}

func (q *QPSClimber) startAdjusting(ctx context.Context) {
	ticker := time.NewTicker(q.options.AdjustInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			q.adjustQPS()
			q.counter.Reset()
		case <-ctx.Done():
			return
		}
	}
}

func (q *QPSClimber) adjustQPS() {
	acceptCount := float64(q.counter.Count(Accept))
	serverOverloadedCount := float64(q.counter.Count(ServerOverloadedCode))
	internalServerErrorCount := float64(q.counter.Count(InternalServerError))
	errorAll := serverOverloadedCount + internalServerErrorCount

	oldRate, newRate := float64(q.limiter.Limit()), float64(0)
	errorRatio := utils.IfGenerics(acceptCount > 0, errorAll/acceptCount, 0)
	switch {
	case errorRatio >= q.options.StrictErrorRatio:
		newRate = oldRate - oldRate*q.options.ProbeRatio
	case errorRatio >= q.options.RelaxErrorRatio:
		newRate = oldRate
	case errorRatio > 0:
		newRate = oldRate + oldRate*q.options.ProbeRatio
	case errorRatio == 0:
		newRate = oldRate + oldRate*q.options.FastProbeRatio

		// 这里防止请求量少的时候，rate无限上涨
		if acceptCount*1.5 < oldRate*q.options.AdjustInterval.Seconds() {
			newRate = oldRate
		}
	}

	finalRate := minMax(1, float64(q.options.MaxRate), newRate)
	if finalRate != oldRate {
		q.limiter.SetLimit(rate.Limit(finalRate))
		q.limiter.SetBurst(int(finalRate))
	}

	fmt.Printf("llm实例：%v 大模型QPS调整: 错误率 %.2f(error=%d, total=%d)，当前limit %d，调整为 %d\n",
		q.uniqueKey, errorRatio, int64(errorAll), int64(acceptCount), int64(oldRate), int64(finalRate))
}

func minMax(min, max, n float64) float64 {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func needRetryError(err error) (statusCode int, retry bool) {
	apiErr := &model.APIError{}
	reqErr := &model.RequestError{}

	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatusCode, apiErr.HTTPStatusCode >= http.StatusInternalServerError || apiErr.HTTPStatusCode == http.StatusTooManyRequests
	} else if errors.As(err, &reqErr) {
		return apiErr.HTTPStatusCode, reqErr.HTTPStatusCode >= http.StatusInternalServerError
	}
	return http.StatusInternalServerError, false
}

type exponentialBackoff struct {
	attempts     int
	initialDelay time.Duration
	maxDelay     time.Duration
	multiplier   float64
}

func newExponentialBackoff(initial, max time.Duration) *exponentialBackoff {
	return &exponentialBackoff{
		initialDelay: initial,
		maxDelay:     max,
		multiplier:   2.0,
	}
}

func (b *exponentialBackoff) wait() {
	delayFloat := float64(b.initialDelay) * math.Pow(b.multiplier, float64(b.attempts))
	delay := time.Duration(delayFloat)

	if delay > b.maxDelay {
		delay = b.maxDelay
	}

	b.attempts++
	time.Sleep(delay)
}

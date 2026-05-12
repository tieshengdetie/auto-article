package qwen

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// Code 定义
type Code int

const (
	SuccessCode Code = iota
	ServerOverloadedCode
	InternalServerError
	UnknownCode
	Accept
	numOfCodes
)

// Counter
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

func (c *counter) Record(code Code)      { c.codes[code].Add(1) }
func (c *counter) Count(code Code) int64 { return c.codes[code].Load() }
func (c *counter) Reset() {
	for i := range c.codes {
		c.codes[i].Store(0)
	}
}

// Options
type Option func(*Options)

func WithAdjustInterval(interval time.Duration) Option {
	return func(o *Options) {
		if interval >= time.Second {
			o.AdjustInterval = interval
		}
	}
}
func WithInitialRate(r int64) Option {
	return func(o *Options) {
		if r > 0 {
			o.InitialRate = r
		}
	}
}
func WithMaxRate(r int64) Option {
	return func(o *Options) {
		if r > 0 {
			o.MaxRate = r
		}
	}
}

type Options struct {
	AdjustInterval   time.Duration
	InitialRate      int64
	MaxRate          int64
	StrictErrorRatio float64
	RelaxErrorRatio  float64
	ProbeRatio       float64
	FastProbeRatio   float64
	RetryBaseDelay   time.Duration
	RetryMaxDelay    time.Duration
}

func defaultOptions() *Options {
	return &Options{
		AdjustInterval:   10 * time.Second,
		InitialRate:      10,
		MaxRate:          1000,
		StrictErrorRatio: 0.05,
		RelaxErrorRatio:  0.01,
		ProbeRatio:       0.05,
		FastProbeRatio:   0.1,
		RetryBaseDelay:   100 * time.Millisecond,
		RetryMaxDelay:    3 * time.Second,
	}
}

// QPSClimber

type QPSClimber struct {
	apiKey  string
	client  *httpClient
	limiter *rate.Limiter
	counter Counter
	options *Options
}

func NewQPSClimber(apiKey string, url string, opts ...Option) *QPSClimber {
	config := defaultOptions()
	for _, opt := range opts {
		opt(config)
	}

	climber := &QPSClimber{
		apiKey:  apiKey,
		client:  newHTTPClient(apiKey, url),
		limiter: rate.NewLimiter(rate.Limit(config.InitialRate), int(config.InitialRate)),
		counter: NewCounter(),
		options: config,
	}

	go climber.startAdjusting(context.Background())
	return climber
}

// DoRequest 调用 OpenAI 风格接口（非流式）
func (q *QPSClimber) DoRequest(ctx context.Context, req *ChatCompletionRequest) (*QwenResponse, error) {

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
		resp, err := q.client.sendQwenRequest(ctx, req)

		// 只有在请求成功且resp不为nil时才记录token使用量
		if err == nil && resp != nil {

			return resp, nil
		}

		if err != nil {
			fmt.Println(fmt.Sprintf("callQwen error: %s", err))
			statusCode, retry := needRetryError(err)
			q.Record(convertToCode(statusCode))

			if !retry {
				return resp, err
			}

			backOff.wait()
		}
	}
}

// DoRequestStream 调用 OpenAI 风格接口（流式）
func (q *QPSClimber) DoRequestStream(ctx context.Context, req *ChatCompletionRequest) (<-chan QwenResponse, error) {

	if !req.Stream {
		req.Stream = true
	}
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
		textCh, errCh := q.client.sendQwenStream(ctx, req)
		var err error
		// 使用select语句同时等待errCh和上下文取消，避免阻塞
		select {
		case <-textCh:

			return textCh, nil
		case err = <-errCh:
			// 打印错误
			fmt.Printf("stream chat error: %v\n", err)

		}

		statusCode, retry := needRetryError(err)
		q.Record(convertToCode(statusCode))

		if !retry {
			return nil, fmt.Errorf("stream request failed after retries")
		}

		backOff.wait()
	}
}

// Record 用于外部记录指标（如 Accept）
func (q *QPSClimber) Record(code Code) { q.counter.Record(code) }

// 启动 QPS 自适应调整
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

// 完善版 QPS 调整策略：根据错误率动态调整
func (q *QPSClimber) adjustQPS() {
	acceptCount := float64(q.counter.Count(Accept))
	if acceptCount == 0 {
		return
	}

	serverOverloadedCount := float64(q.counter.Count(ServerOverloadedCode))
	internalServerErrorCount := float64(q.counter.Count(InternalServerError))
	errorAll := serverOverloadedCount + internalServerErrorCount

	oldRate := float64(q.limiter.Limit())
	newRate := float64(0)

	errorRatio := 0.0
	if acceptCount > 0 {
		errorRatio = errorAll / acceptCount
	}

	switch {
	case errorRatio >= q.options.StrictErrorRatio:
		// 错误率高于严格阈值，降低 QPS
		newRate = oldRate - oldRate*q.options.ProbeRatio
	case errorRatio >= q.options.RelaxErrorRatio:
		// 错误率在宽松阈值和严格阈值之间，保持 QPS 不变
		newRate = oldRate
	case errorRatio > 0:
		// 有少量错误，缓慢增加 QPS
		newRate = oldRate + oldRate*q.options.ProbeRatio
	case errorRatio == 0:
		// 没有错误，快速增加 QPS
		newRate = oldRate + oldRate*q.options.FastProbeRatio

		// 防止请求量少时 QPS 无限上涨
		if acceptCount*1.5 < oldRate*q.options.AdjustInterval.Seconds() {
			newRate = oldRate
		}
	}

	// 确保 QPS 在合理范围内
	finalRate := math.Max(1, math.Min(float64(q.options.MaxRate), newRate))
	if finalRate != oldRate {
		q.limiter.SetLimit(rate.Limit(finalRate))
		q.limiter.SetBurst(int(finalRate))
	}
}

// 辅助函数：转换HTTP状态码为Code类型
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

// 辅助函数：判断是否需要重试
func needRetryError(err error) (statusCode int, retry bool) {
	// 从错误字符串中解析HTTP状态码
	// 错误格式通常为："qwen api error: status=429, body=..." 或 "stream error: status=429, body=..."
	errMsg := err.Error()
	statusPrefix := "status="
	statusStart := strings.Index(errMsg, statusPrefix)
	if statusStart != -1 {
		statusStart += len(statusPrefix)
		statusEnd := strings.Index(errMsg[statusStart:], ",")
		if statusEnd != -1 {
			statusCodeStr := errMsg[statusStart : statusStart+statusEnd]
			if code, err := strconv.Atoi(statusCodeStr); err == nil {
				statusCode = code
			} else {
				statusCode = http.StatusInternalServerError
			}
		} else {
			statusCode = http.StatusInternalServerError
		}
	} else {
		statusCode = http.StatusInternalServerError
	}

	// 判断是否需要重试
	retry = statusCode >= http.StatusInternalServerError || statusCode == http.StatusTooManyRequests
	return
}

// 指数退避重试结构体
type exponentialBackoff struct {
	attempts     int
	initialDelay time.Duration
	maxDelay     time.Duration
	multiplier   float64
}

// 创建指数退避重试实例
func newExponentialBackoff(initial, max time.Duration) *exponentialBackoff {
	return &exponentialBackoff{
		initialDelay: initial,
		maxDelay:     max,
		multiplier:   2.0,
	}
}

// 等待重试
func (b *exponentialBackoff) wait() {
	delayFloat := float64(b.initialDelay) * math.Pow(b.multiplier, float64(b.attempts))
	delay := time.Duration(delayFloat)

	if delay > b.maxDelay {
		delay = b.maxDelay
	}

	b.attempts++
	time.Sleep(delay)
}

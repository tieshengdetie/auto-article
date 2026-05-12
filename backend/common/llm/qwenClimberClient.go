package llm

import (
	"AutoArticle/common/climber/qwen"
	"AutoArticle/common/llm/types"
	"AutoArticle/global"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// QwenOption 定义调用Qwen大模型的选项
type QwenOption struct {
	ModelName     string   // 模型名称
	SystemMessage []string // 系统消息
	UserMessage   string   // 用户消息
	ImgUrl        []string // 图片URL列表
	Temperature   float64  // 温度参数
}

type QwenClimberClient struct {
	QpsClimber *qwen.QPSClimber
	ApiKey     string
}

const (
	QwenApiUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
)

func NewQwenClimberClient(apiKey string) *QwenClimberClient {

	qpsClimber := qwen.NewQPSClimber(apiKey, QwenApiUrl,
		// qps速率限制更新间隔，默认10s调整一次
		qwen.WithAdjustInterval(10*time.Second),
		// 初始qps速率，默认为10
		qwen.WithInitialRate(50),
		// qps速率最大限制，默认为1000
		// qps速率更新到该值后不再升高，请根据业务体量自行调整
		qwen.WithMaxRate(1000))
	return &QwenClimberClient{
		QpsClimber: qpsClimber,
		ApiKey:     apiKey,
	}
}

// CallQwen 调用Qwen大模型（非流式）
func (q *QwenClimberClient) CallQwen(ctx context.Context, option types.QwenOption) ([]qwen.QwenChoice, int, error) {

	var resp *qwen.QwenResponse
	var err error
	// 记录开始时间
	start := time.Now()
	// 调用Qwen大模型
	if option.Stream {
		resp, err = q.DoRequestWithStream(ctx, option)
	} else {
		resp, err = q.DoRequest(ctx, option)
	}

	// 计算耗时
	duration := time.Since(start)
	seconds := duration.Seconds()
	// 错误处理
	if err != nil {
		global.Logger.Sugar().Errorf("大模型请求失败:resp:%v", resp)
		global.Logger.Sugar().Errorf("大模型请求失败:err:%v", err)
		// 根据错误类型返回相应的错误码
		if strings.Contains(err.Error(), `"type":"Unauthorized"`) || strings.Contains(err.Error(), `Error code: 401`) {
			return nil, 401, errors.New("API Key 或 AK/SK 校验未通过")
		}
		if strings.Contains(err.Error(), `"type":"TooManyRequests"`) || strings.Contains(err.Error(), `Error code: 429`) {
			return nil, 429, errors.New("请求过于频繁")
		}
		return nil, 0, errors.Wrap(err, "大模型请求失败")
	}

	// 长时间请求日志记录
	if seconds > 60 {
		respStr := ""
		if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
			respStr = resp.Choices[0].Message.Content
		}

		global.Logger.Sugar().Infof("千问耗时: %s %.2f秒 模型:%s %s", resp.ID, seconds, option.ModelName, respStr)
	}
	if resp != nil && len(resp.Choices) > 0 && resp.Choices[0].Message.Content == "" {
		global.Logger.Sugar().Infof("千问耗时: %s %.2f秒 模型:%s 千问返回结果为空", resp.ID, seconds, option.ModelName)

	}
	// 返回结果
	return resp.Choices, http.StatusOK, nil
}

func (q *QwenClimberClient) DoRequest(ctx context.Context, option types.QwenOption) (*qwen.QwenResponse, error) {
	req := q.getRequest(option)
	qwenResp, err := q.QpsClimber.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return qwenResp, nil
}

func (q *QwenClimberClient) getRequest(option types.QwenOption) *qwen.ChatCompletionRequest {
	// 构建消息列表
	messages := make([]qwen.QwenMessage, 0)

	// 添加系统消息
	for _, system := range option.SystemMessage {
		if system != "" {
			messages = append(messages, qwen.QwenMessage{
				Role:    "system",
				Content: system,
			})
		}
	}

	// 根据是否有图片决定构建不同类型的用户消息
	if len(option.ImgUrl) > 0 {
		// 多模态调用，支持图片
		// 使用map切片来构建多模态内容
		userValues := make([]map[string]interface{}, 0)

		// 添加文本部分（如果不为空）
		if option.UserMessage != "" {
			userValues = append(userValues, map[string]interface{}{
				"type": "text",
				"text": option.UserMessage,
			})
		}

		// 添加图片URL
		for _, url := range option.ImgUrl {
			userValues = append(userValues, map[string]interface{}{
				"type":      "image_url",
				"image_url": url,
			})
		}

		// 添加用户消息（包含文本和图片）
		messages = append(messages, qwen.QwenMessage{
			Role:    "user",
			Content: userValues,
		})

	} else {
		// 纯文本调用
		if option.UserMessage != "" {
			messages = append(messages, qwen.QwenMessage{
				Role:    "user",
				Content: option.UserMessage,
			})
		}
	}
	req := &qwen.ChatCompletionRequest{
		Model:    option.ModelName,
		Messages: messages,
	}

	// 设置默认温度
	temperature := option.Temperature
	if temperature == 0 {
		temperature = 0.7
	}
	req.Temperature = &temperature
	if option.Stream {
		req.Stream = true
	}
	// 设置深度思考和最大思考长度
	req.EnableThinking = option.EnableThinking
	req.ThinkingBudget = 10000
	if option.ThinkingBudget > 0 {
		req.ThinkingBudget = option.ThinkingBudget
	}

	if option.MaxTokens > 0 {
		req.MaxTokens = &option.MaxTokens
	}

	return req
}

// DoRequestWithStream 调用Qwen大模型（流式）
func (q *QwenClimberClient) DoRequestWithStream(ctx context.Context, option types.QwenOption) (*qwen.QwenResponse, error) {
	req := q.getRequest(option)
	// 发送流式请求
	stream, err := q.QpsClimber.DoRequestStream(ctx, req)
	if err != nil {
		return nil, err
	}

	// 读取所有流式响应
	var fullContent string
	var responseID string
	var finishReason, role, object string
	var usage qwen.QwenUsage

	// 处理每个响应块
	for chunk := range stream {
		// 保存响应ID
		if chunk.ID != "" {
			responseID = chunk.ID
		}

		// 处理每个选择
		for _, choice := range chunk.Choices {
			// 处理内容
			if choice.Delta != nil {
				fullContent += choice.Delta.Content
			}
			if choice.Delta.Role != "" {
				role = choice.Delta.Role
			}

			// 处理完成原因
			if choice.FinishReason != "" {
				finishReason = choice.FinishReason
			}
		}

		// 保存Usage信息（如果有）
		if chunk.Usage.TotalTokens > 0 {
			usage = chunk.Usage
		}
		if chunk.Object != "" {
			object = chunk.Object
		}
		//fmt.Println("fullContent:", fullContent)
	}

	// 如果没有获取到响应
	if fullContent == "" {
		return nil, errors.New("未获取到流式响应内容")
	}

	// 创建最终的响应结构
	finalResponse := &qwen.QwenResponse{
		ID:      responseID,
		Object:  object,
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []qwen.QwenChoice{
			{
				Index: 0,
				Message: qwen.QwenRespMessage{
					Role:    role,
					Content: fullContent,
				},
				FinishReason: finishReason,
			},
		},
		Usage: usage,
	}

	return finalResponse, nil
}

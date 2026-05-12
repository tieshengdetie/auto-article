package llm

import (
	"AutoArticle/common/climber"
	"AutoArticle/common/llm/types"
	"AutoArticle/global"
	"AutoArticle/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

// ModelMap 模型map映射，模型的enpId更改时要对应更改

// SchemaGenerator 是一个泛型函数类型，用于生成JSON Schema

type DBaoClimberClient struct {
	QpsClimber    *climber.QPSClimber
	LLMGateWayUrl string
}

// 默认的Schema生成函数

func GenerateSchema[T any]() *jsonschema.Schema { // <-- 优化返回类型为具体 Schema 类型
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	return reflector.Reflect(new(T)) // 使用 new(T) 避免空值问题
}

func NewDBaoClimberClient(apiKey string, llmBaseUrl string) *DBaoClimberClient {
	// 请确保您已将 API Key 存储在环境变量 ARK_API_KEY 中
	// 初始化Ark客户端，从环境变量中读取您的API Key
	client := arkruntime.NewClientWithApiKey(
		// 从环境变量中获取您的 API Key。此为默认方式，您可根据需要进行修改
		apiKey,
		// 此为默认路径，您可根据业务所在地域进行配置
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
		arkruntime.WithRetryTimes(0),
	)
	qpsClimber := climber.NewQPSClimber(client, apiKey,
		// qps速率限制更新间隔，默认10s调整一次
		climber.WithAdjustInterval(10*time.Second),
		// 初始qps速率，默认为10
		climber.WithInitialRate(50),
		// qps速率最大限制，默认为1000
		// qps速率更新到该值后不再升高，请根据业务体量自行调整
		climber.WithMaxRate(1000))
	return &DBaoClimberClient{
		QpsClimber:    qpsClimber,
		LLMGateWayUrl: llmBaseUrl,
	}
}

/************************************************大模型结构化输出方法**********************************************************************/

// CallDBao
//
//	@Description:
//	@receiver d
//	@param ctx
//	@param option 参数说明:
//	ImgUrl        []string  图片地址
//	SystemMessage []string  系统消息
//	UserMessage   string  用户消息
//	ModelName     string  模型名称
//	Temperature   float32  温度
//	Thinking      string  是否开启思考模式
//	ThinkingLevel string  思考级别
//	Schema  	  *jsonschema.Schema 结构化结构体
//
// @return []*model.ChatCompletionChoice
// @return int
// @return error
func (d *DBaoClimberClient) CallDBao(ctx context.Context, option types.DBaoOption) (*model.ChatCompletionResponse, int, error) {

	// 记录开始时间
	start := time.Now()
	var resp model.ChatCompletionResponse
	var err error
	// 发送请求
	if option.Stream == true {
		resp, err = d.DoRequestWithStream(ctx, option)
	} else {
		resp, err = d.DoRequest(ctx, option)
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
		return nil, 0, errors.New("请求第三方错误")
	}

	// 长时间请求日志记录
	if seconds > 60 {
		if len(resp.Choices) > 0 && resp.Choices[0].Message.Content.StringValue != nil {

		}

		global.Logger.Sugar().Infof("豆包耗时: %s %.2f秒 模型:%s", resp.ID, seconds, option.ModelName)
	}
	//respStr := *resp.Choices[0].Message.Content.StringValue
	//fmt.Println(len(respStr))
	return &resp, http.StatusOK, nil
}

// DoRequestWithStream
//
//	@Description: 流式调用大模型，但返回与CallDBao相同的结果格式
//	@receiver d
//	@param ctx
//	@param req 请求参数
//
// @return model.ChatCompletionResponse 与CallDBao相同的结果格式
// @return error 错误信息
func (d *DBaoClimberClient) DoRequestWithStream(ctx context.Context, option types.DBaoOption) (model.ChatCompletionResponse, error) {
	req := d.getRequest(option)
	// 发送流式请求
	streamReader, err := d.QpsClimber.DoRequestStream(ctx, req)
	if err != nil {
		global.Logger.Sugar().Errorf("大模型流式请求失败:err:%v", err)
		return model.ChatCompletionResponse{}, err
	}

	// 读取所有流式响应
	var fullContent string
	var finishReason model.FinishReason
	var role string
	var responseID string

	for {
		chunk, err := streamReader.Recv()
		if err != nil {
			// 流结束
			break
		}

		// 保存响应ID
		if chunk.ID != "" {
			responseID = chunk.ID
		}

		// 处理每个选择
		for _, choice := range chunk.Choices {
			// 设置角色（通常在第一个chunk中）
			if role == "" {
				role = choice.Delta.Role
			}

			// 处理内容
			if choice.Delta.Content != "" {
				fullContent += choice.Delta.Content
			}

			// 处理完成原因
			if choice.FinishReason != "" {
				finishReason = choice.FinishReason
			}
		}
		//fmt.Println("fullContent:", fullContent)
	}

	// 关闭流
	_ = streamReader.Close()

	// 如果没有获取到响应
	if fullContent == "" {
		return model.ChatCompletionResponse{}, errors.New("未获取到流式响应内容")
	}

	// 创建最终的响应结构，与DoRequest返回格式一致
	finalResponse := model.ChatCompletionResponse{
		ID:      responseID, // 使用最后一个chunk的ID
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []*model.ChatCompletionChoice{
			{
				Index: 0,
				Message: model.ChatCompletionMessage{
					Role: role,
					Content: &model.ChatCompletionMessageContent{
						StringValue: &fullContent,
					},
				},
				FinishReason: finishReason,
			},
		},
	}

	return finalResponse, nil
}

func (d *DBaoClimberClient) DoRequest(ctx context.Context, option types.DBaoOption) (model.ChatCompletionResponse, error) {

	req := d.getRequest(option)
	resp, err := d.QpsClimber.DoRequest(ctx, req)
	return resp, err

}

func (d *DBaoClimberClient) getRequest(option types.DBaoOption) *model.CreateChatCompletionRequest {
	// 构建消息列表
	messages := make([]*model.ChatCompletionMessage, 0)

	// 添加系统消息
	for _, system := range option.SystemMessage {
		if system != "" {
			messages = append(messages, &model.ChatCompletionMessage{
				Role:    model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{StringValue: volcengine.String(system)},
			})
		}
	}

	// 构建请求参数
	temperature := option.Temperature
	if temperature == 0 {
		temperature = 0.7 // 设置默认温度
	}

	// 根据是否有图片决定构建不同类型的用户消息
	if len(option.ImgUrl) > 0 {
		// 多模态调用，支持图片
		userValues := make([]*model.ChatCompletionMessageContentPart, 0)

		// 添加文本部分（如果不为空）
		if option.UserMessage != "" {
			userValues = append(userValues, &model.ChatCompletionMessageContentPart{
				Type: model.ChatCompletionMessageContentPartTypeText,
				Text: option.UserMessage,
			})
		}

		// 添加图片URL，处理URL替换
		for _, url := range option.ImgUrl {
			// 替换图片URL域名
			url = strings.ReplaceAll(url, "https://homework-img.readboy.com", "https://readboy-homework.oss-cn-shenzhen.aliyuncs.com")
			imageDetail := model.ImageURLDetailHigh
			if option.ImageURLDetail == "low" {
				imageDetail = model.ImageURLDetailLow
			} else if option.ImageURLDetail == "auto" {
				imageDetail = model.ImageURLDetailAuto
			}
			userValues = append(userValues, &model.ChatCompletionMessageContentPart{
				Type:     model.ChatCompletionMessageContentPartTypeImageURL,
				ImageURL: &model.ChatMessageImageURL{URL: url, Detail: imageDetail},
			})
		}

		// 只有当userValues不为空时才添加用户消息（至少包含图片）
		messages = append(messages, &model.ChatCompletionMessage{
			Role: model.ChatMessageRoleUser,
			Content: &model.ChatCompletionMessageContent{
				ListValue: userValues,
			},
		})

	} else {
		// 纯文本调用
		if option.UserMessage != "" {
			messages = append(messages, &model.ChatCompletionMessage{
				Role:    model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{StringValue: volcengine.String(option.UserMessage)},
			})
		}
	}
	// 此处兼容以前的方式
	modelName := utils.IfGenerics(option.EndPointId != "", option.EndPointId, option.ModelName)
	// 构建请求
	req := &model.CreateChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		Temperature: &temperature,
	}
	// 限制输出长度
	if option.MaxCompletionTokens > 0 {
		req.MaxCompletionTokens = &option.MaxCompletionTokens
	}
	// 设置思考参数
	if option.Thinking != "" {
		thinkingModel := model.Thinking{Type: model.ThinkingTypeAuto}
		if option.Thinking == "true" {
			thinkingModel = model.Thinking{Type: model.ThinkingTypeEnabled}
		} else if option.Thinking == "false" {
			thinkingModel = model.Thinking{Type: model.ThinkingTypeDisabled}
		}
		req.Thinking = &thinkingModel

		// 设置思考级别
		if option.ThinkingLevel != "" && option.Thinking == "true" {
			reasoningEffort := model.ReasoningEffortLow
			switch option.ThinkingLevel {
			case "minimal":
				reasoningEffort = model.ReasoningEffortMinimal
			case "low":
				reasoningEffort = model.ReasoningEffortLow
			case "medium":
				reasoningEffort = model.ReasoningEffortMedium
			case "high":
				reasoningEffort = model.ReasoningEffortHigh
			default:
				reasoningEffort = model.ReasoningEffortLow
			}
			req.ReasoningEffort = &reasoningEffort
		}
	}

	// 设置结构化输出（如果提供了Schema）
	if option.Schema != nil {
		// 启用JSON格式输出
		responseFormat := model.ResponseFormat{
			Type: model.ResponseFormatJsonObject,
		}
		req.ResponseFormat = &responseFormat
	}
	if option.MaxTokens > 0 {
		req.MaxTokens = volcengine.Int(option.MaxTokens)
	}
	if option.Stream {
		req.Stream = volcengine.Bool(true)
	}

	return req

}

package sendHttp

import (
	"AutoArticle/common/sendHttp/types"
	"AutoArticle/global"
	"AutoArticle/utils"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

const AliYunModelType = "aliYunModel"

var AliYunModel = AliYunModelHttp{}

type AliYunModelHttp struct{}

var AliYunModelHttpUrlMap = urlMap{
	"askQuestion": "/v1/chat/completions",
}

func init() {
	registerUrl(AliYunModelType, AliYunModelHttpUrlMap)
}
func (a *AliYunModelHttp) Config(ctx context.Context, urlKey string) (headers map[string]string, fullUrl string) {
	//设置请求头
	headers = make(map[string]string)
	config := global.ServerConfig.AliYunModel
	headers["Authorization"] = "Bearer " + config.AppKey

	fullUrl = config.BaseUrl + urlMapping[AliYunModelType][urlKey]
	return
}
func (a *AliYunModelHttp) Request(ctx context.Context, method int, data map[string]interface{}, urlKey string, receiver interface{}) (err error) {
	headers, fullUrl := a.Config(ctx, urlKey)
	return utils.HTTPRequestByReceiver(method, fullUrl, data, headers, nil, receiver)
}

// AskQuestion
//
//	@Description: 向模型提问
//	@receiver a
//	@param ctx
//	@param message
func (a *AliYunModelHttp) AskQuestion(ctx context.Context, message types.Message) (resp map[string]interface{}, err error) {

	var messages = make([]types.Message, 0)
	messages = append(messages, types.Message{Role: "system", Content: "你是一个AI助手，请根据用户的问题给出简洁明了的答案"})
	messages = append(messages, message)
	var param = map[string]interface{}{
		"model":    "qwen-plus",
		"messages": messages,
	}
	var pointResp = map[string]interface{}{}
	err = a.Request(ctx, utils.PostJson, param, "askQuestion", &pointResp)
	if err != nil {
		return
	}
	return pointResp, nil
}

// ImageRecognition
//
//	@Description: 图像识别
//	@receiver a
//	@param ctx
//	@param message
//	@return vo
//	@return err
func (a *AliYunModelHttp) ImageRecognition(ctx context.Context, message types.Message) (resp map[string]interface{}, err error) {

	var messages = make([]types.Message, 0)
	messages = append(messages, message)
	var param = map[string]interface{}{
		"model":    "qwen-vl-max-latest",
		"messages": messages,
	}
	var pointResp = map[string]interface{}{}
	err = a.Request(ctx, utils.PostJson, param, "askQuestion", &pointResp)
	if err != nil {
		return
	}
	return pointResp, nil
}

// AskQuestionByStream
//
//	@Description: 流式提问
//	@receiver a
//	@param ctx
//	@param message
//	@return vo
//	@return err
func (a *AliYunModelHttp) AskQuestionByStream(ctx context.Context, message types.Message) (<-chan string, <-chan error) {
	//设置请求头
	headers := make(map[string]string)
	config := global.ServerConfig.AliYunModel
	headers["Authorization"] = "Bearer " + config.AppKey

	fullUrl := config.BaseUrl + urlMapping[AliYunModelType]["askQuestion"]

	lines := make(chan string)
	errors := make(chan error, 1)
	var messages = make([]types.Message, 0)
	messages = append(messages, types.Message{Role: "system", Content: "你是一个AI助手，请根据用户的问题给出简洁明了的答案"})
	messages = append(messages, message)
	var param = map[string]interface{}{
		"model":    "qwen-plus",
		"messages": messages,
		"stream":   true,
	}
	go func() {
		defer close(lines)
		defer close(errors)

		stream, err := utils.HTTPRequestByStream(utils.PostJson, fullUrl, param, headers, nil)
		if err != nil {
			return
		}
		defer stream.Body.Close()

		scanner := bufio.NewScanner(stream.Body)
		for scanner.Scan() {
			// 解析每一行数据
			parseLine := streamDataParse(scanner.Bytes())
			fmt.Println(string(parseLine))
			if len(parseLine) == 0 {

				continue
			}

			// 检查是否以 "data:" 开头
			if string(parseLine) == "sse-invalid-data-flag" {
				// 协议头，不用处理
				continue
			}
			var answer types.ChatCompletionChunk
			if err := json.Unmarshal(parseLine, &answer); err != nil {
				return
			}
			lines <- answer.Choices[0].Delta.Content
		}

		if err := scanner.Err(); err != nil {
			errors <- fmt.Errorf("error reading response: %v", err)
		}
	}()

	return lines, errors
}

// streamDataParse 流式输出处理
func streamDataParse(line []byte) []byte {
	// 可能返回空格字符串
	trimMsg := bytes.TrimSpace(line)

	if len(trimMsg) == 0 {
		return []byte{}
	}

	// 检查是否以 "data:" 开头
	if !strings.HasPrefix(string(trimMsg), "data:") {
		return []byte("sse-invalid-data-flag")
	}

	// 接收处理数据
	trimmedLine := strings.TrimPrefix(string(trimMsg), "data:")

	return []byte(trimmedLine + "\n")
}

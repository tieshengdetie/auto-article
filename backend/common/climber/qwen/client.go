package qwen

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChatCompletionRequest struct {
	Model          string        `json:"model"`
	Messages       []QwenMessage `json:"messages"`
	MaxTokens      *int          `json:"max_tokens,omitempty"`
	Temperature    *float64      `json:"temperature"`
	TopP           *float64      `json:"top_p,omitempty"`
	Stream         bool          `json:"stream"`
	EnableThinking bool          `json:"enable_thinking"`
	ThinkingBudget int           `json:"thinking_budget"`
}

type QwenMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // 支持字符串或多模态内容数组
}

type qwenChatRequest struct {
	Model      string     `json:"model"`
	Input      qwenInput  `json:"input"`
	Parameters qwenParams `json:"parameters,omitempty"`
}

type qwenInput struct {
	Messages []QwenMessage `json:"messages"`
}

type qwenParams struct {
	ResultFormat string  `json:"result_format,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	TopP         float64 `json:"top_p,omitempty"`
	Stream       bool    `json:"stream,omitempty"`
}

type QwenResponse struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"`
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Choices           []QwenChoice `json:"choices"`
	Usage             QwenUsage    `json:"usage"`
	SystemFingerprint *string      `json:"system_fingerprint,omitempty"`
	Code              int          `json:"code,omitempty"`
	Message           string       `json:"message,omitempty"`
}

type QwenChoice struct {
	Index        int             `json:"index"`
	Message      QwenRespMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Logprobs     *interface{}    `json:"logprobs,omitempty"`
	Delta        *Delta          `json:"delta"`
}
type Delta struct {
	Content      string      `json:"content"`
	FunctionCall interface{} `json:"function_call"` // 可能为 null 或对象，若结构已知可替换为具体类型
	Refusal      interface{} `json:"refusal"`       // 同上，可能为 null 或字符串等
	Role         string      `json:"role"`
	ToolCalls    interface{} `json:"tool_calls"` // 可能为 null 或数组，若结构已知可定义为 []ToolCall
}
type QwenRespMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type qwenOutput struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type QwenUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	PromptTokensDetails     map[string]int `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails map[string]int `json:"completion_tokens_details,omitempty"`
}

type qwenStreamEvent struct {
	Output  qwenStreamOutput `json:"output"`
	Usage   QwenUsage        `json:"usage,omitempty"`
	Code    int              `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
}

type qwenStreamOutput struct {
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// httpClient 封装原始 HTTP 调用
type httpClient struct {
	apiKey string
	url    string
	client *http.Client
}

func newHTTPClient(apiKey, url string) *httpClient {
	return &httpClient{
		apiKey: apiKey,
		url:    url,
		client: &http.Client{Timeout: 900 * time.Second},
	}
}

// sendQwenRequest 发送非流式请求
func (h *httpClient) sendQwenRequest(ctx context.Context, req *ChatCompletionRequest) (*QwenResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	//关闭绿网防止400：X-DashScope-DataInspection：{"input":"disable","output":"disable"}
	httpReq.Header.Set("X-DashScope-DataInspection", "{\"input\":\"disable\",\"output\":\"disable\"}")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qwen api error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	//fmt.Println(string(respBody))
	var qwenResp QwenResponse
	if err := json.Unmarshal(respBody, &qwenResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if qwenResp.Code != 0 && qwenResp.Code != 200 {
		return nil, fmt.Errorf("qwen logic error: code=%d, msg=%s", qwenResp.Code, qwenResp.Message)
	}

	return &qwenResp, nil
}

// sendQwenStream 发送流式请求
func (h *httpClient) sendQwenStream(ctx context.Context, req *ChatCompletionRequest) (<-chan QwenResponse, <-chan error) {
	textCh := make(chan QwenResponse, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(textCh)
		defer close(errCh)

		body, err := json.Marshal(req)
		if err != nil {
			errCh <- fmt.Errorf("marshal stream request: %w", err)
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, bytes.NewBuffer(body))
		if err != nil {
			errCh <- fmt.Errorf("create stream http request: %w", err)
			return
		}
		httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "text/event-stream")

		resp, err := h.client.Do(httpReq)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			errCh <- fmt.Errorf("stream error: status=%d, body=%s", resp.StatusCode, string(respBody))
			return
		}

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
				return
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				if data == "[DONE]" {
					break
				}

				var event QwenResponse
				if err := json.Unmarshal([]byte(data), &event); err != nil {
					continue // skip malformed
				}

				if event.Code != 0 && event.Code != 200 {
					errCh <- fmt.Errorf("stream api error: code=%d, msg=%s", event.Code, event.Message)
					return
				}

				select {
				case textCh <- event:
				case <-ctx.Done():
					return
				}

			}
		}
	}()

	return textCh, errCh
}

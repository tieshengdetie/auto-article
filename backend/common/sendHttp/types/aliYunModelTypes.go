package types

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// 定义 Delta 结构体

type Delta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// 定义 Choice 结构体

type Choice struct {
	Delta        Delta   `json:"delta"`
	Index        int     `json:"index"`
	Logprobs     *int    `json:"logprobs"`      // 使用指针便于处理 null 值
	FinishReason *string `json:"finish_reason"` // 使用指针便于处理 null 值
}

// 定义主响应结构体

type ChatCompletionChunk struct {
	Choices           []Choice `json:"choices"`
	Object            string   `json:"object"`
	Usage             *string  `json:"usage"` // 使用指针便于处理 null 值
	Created           int64    `json:"created"`
	SystemFingerprint *string  `json:"system_fingerprint"` // 使用指针便于处理 null 值
	Model             string   `json:"model"`
	ID                string   `json:"id"`
}

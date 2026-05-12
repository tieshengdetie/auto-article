package types

import "github.com/google/gnostic-models/jsonschema"

type DBaoOption struct {
	ImgUrl              []string
	SystemMessage       []string
	UserMessage         string
	ModelName           string
	ModelNameNew        string
	EndPointId          string
	Temperature         float32
	Thinking            string
	ThinkingLevel       string
	ImageURLDetail      string // "high"、"low"、"auto"、空字符串时为"high"
	Schema              *jsonschema.Schema
	MaxTokens           int
	MaxCompletionTokens int
	Stream              bool
}
type QwenOption struct {
	ModelName      string   // 模型名称
	SystemMessage  []string // 系统消息
	UserMessage    string   // 用户消息
	ImgUrl         []string // 图片URL列表
	Temperature    float64  // 温度参数
	Stream         bool
	MaxTokens      int
	EnableThinking bool
	ThinkingBudget int // 最大思考长度，默认设置15000
}

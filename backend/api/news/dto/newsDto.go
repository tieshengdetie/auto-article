package dto

// AllNewsReq 指定分类新闻请求参数
type AllNewsReq struct {
	Col  int    `json:"col" validate:"required,min=1"`           // 频道ID（必填）
	Num  int    `json:"num" validate:"min=1,max=50"`             // 返回数量（1-50，默认10）
	Page int    `json:"page" validate:"min=0"`                   // 页码（默认1）
	Rand int    `json:"rand" validate:"oneof=0 1"`               // 是否随机获取（0不随机，1随机）
	Word string `json:"word"`                                    // 搜索关键词（可选）
}

// GeneralNewsReq 综合新闻请求参数
type GeneralNewsReq struct {
	Num    int    `json:"num" validate:"min=1,max=50"`           // 返回数量（1-50，默认10）
	Page   int    `json:"page" validate:"min=0"`                   // 页码（默认0）
	Rand   int    `json:"rand" validate:"oneof=0 1"`               // 是否随机获取（0不随机，1随机）
	Word   string `json:"word"`                                    // 搜索关键词（可选）
	Source string `json:"source"`                                  // 指定来源（可选）
}

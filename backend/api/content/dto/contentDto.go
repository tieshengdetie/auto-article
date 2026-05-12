package dto

type NewsItemReq struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Ctime       string `json:"ctime"`
	Title       string `json:"title"`
	PicURL      string `json:"picUrl"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

type SaveNewsReq struct {
	Keyword     string        `json:"keyword"`
	SearchType  string        `json:"searchType"`
	ChannelID   int           `json:"channelId"`
	ChannelName string        `json:"channelName"`
	List        []NewsItemReq `json:"list"`
}

type ListNewsReq struct {
	Keyword    string `form:"keyword"`
	SearchType string `form:"searchType"`
	ChannelID  int    `form:"channelId"`
}

type ListArticlesReq struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type GenerateArticleReq struct {
	Keyword     string `json:"keyword"`
	UserPrompt  string `json:"userPrompt"`
	LLMProvider string `json:"llmProvider"`
	ModelName   string `json:"modelName"`
	NewsIDs     []uint `json:"newsIds"`
}

type UpdateArticleReq struct {
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent"`
	CoverImageURL   string `json:"coverImageUrl"`
	Tags            string `json:"tags"`
}

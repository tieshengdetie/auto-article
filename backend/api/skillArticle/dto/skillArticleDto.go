package dto

type CreateSkillArticleReq struct {
	TaskID          string `json:"taskId"`
	Platform        string `json:"platform" binding:"required"`
	Keyword         string `json:"keyword" binding:"required"`
	KeywordSegments string `json:"keywordSegments"`
	Category        string `json:"category"`
	Title           string `json:"title" binding:"required"`
	TitleOptions    string `json:"titleOptions"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent" binding:"required"`
	Tags            string `json:"tags"`
	CoverImageURL   string `json:"coverImageUrl"`
	CoverImageType  string `json:"coverImageType"`
	Images          string `json:"images"`
	Sources         string `json:"sources"`
	HotTopics       string `json:"hotTopics"`
	StyleProfile    string `json:"styleProfile"`
	WordCount       int    `json:"wordCount"`
	ModelProvider   string `json:"modelProvider"`
	ModelName       string `json:"modelName"`
	PromptVersion   string `json:"promptVersion"`
	SkillVersion    string `json:"skillVersion"`
	HumanizeStatus  string `json:"humanizeStatus"`
	Status          string `json:"status"`
	PublishStatus   string `json:"publishStatus"`
	PublishPayload  string `json:"publishPayload"`
	ErrorMessage    string `json:"errorMessage"`
}

type ListSkillArticlesReq struct {
	Keyword  string `form:"keyword"`
	Platform string `form:"platform"`
	Category string `form:"category"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type UpdateSkillArticleReq struct {
	Title           string `json:"title"`
	TitleOptions    string `json:"titleOptions"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent"`
	Tags            string `json:"tags"`
	CoverImageURL   string `json:"coverImageUrl"`
	CoverImageType  string `json:"coverImageType"`
	Images          string `json:"images"`
	Sources         string `json:"sources"`
	StyleProfile    string `json:"styleProfile"`
	WordCount       int    `json:"wordCount"`
	HumanizeStatus  string `json:"humanizeStatus"`
	Status          string `json:"status"`
	PublishStatus   string `json:"publishStatus"`
	PublishPayload  string `json:"publishPayload"`
	ErrorMessage    string `json:"errorMessage"`
}

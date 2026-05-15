package dto

type CreateSkillArticleReq struct {
	TaskID          string `json:"taskId"`
	Platform        string `json:"platform" binding:"required"`
	Keyword         string `json:"keyword" binding:"required"`
	Category        string `json:"category"`
	Title           string `json:"title" binding:"required"`
	TitleOptions    string `json:"titleOptions"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent" binding:"required"`
	CoverImageURL   string `json:"coverImageUrl"`
}

type ListSkillArticlesReq struct {
	Keyword  string `form:"keyword"`
	Platform string `form:"platform"`
	Category string `form:"category"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type UpdateSkillArticleReq struct {
	Title           string `json:"title"`
	TitleOptions    string `json:"titleOptions"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent"`
	CoverImageURL   string `json:"coverImageUrl"`
}

package vo

import "AutoArticle/model"

type SkillArticlePageResp struct {
	List     []model.SkillGeneratedArticle `json:"list"`
	Total    int64                         `json:"total"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"pageSize"`
}

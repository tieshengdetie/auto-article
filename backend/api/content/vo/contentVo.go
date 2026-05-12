package vo

import "AutoArticle/model"

type SaveNewsResp struct {
	Saved   int                   `json:"saved"`
	Skipped int                   `json:"skipped"`
	List    []model.CollectedNews `json:"list"`
}

type PublishPackageResp struct {
	Package model.PublishPackage   `json:"package"`
	Article model.GeneratedArticle `json:"article"`
}

type ArticlePageResp struct {
	List     []model.GeneratedArticle `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}

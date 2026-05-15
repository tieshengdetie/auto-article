package service

import (
	"AutoArticle/api/skillArticle/dto"
	"AutoArticle/api/skillArticle/vo"
	"AutoArticle/global"
	"AutoArticle/model"
	"AutoArticle/utils"
	"context"
	"errors"
	"strings"
)

type SkillArticleService struct{}

func (s *SkillArticleService) Create(ctx context.Context, req dto.CreateSkillArticleReq) (*model.SkillGeneratedArticle, error) {
	if strings.TrimSpace(req.Platform) == "" {
		return nil, errors.New("平台不能为空")
	}
	if strings.TrimSpace(req.Keyword) == "" {
		return nil, errors.New("关键词不能为空")
	}
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.MarkdownContent) == "" {
		return nil, errors.New("标题和正文不能为空")
	}
	article := model.SkillGeneratedArticle{
		TaskID:          normalizeTaskID(req.TaskID),
		Platform:        strings.TrimSpace(req.Platform),
		Keyword:         strings.TrimSpace(req.Keyword),
		Category:        req.Category,
		Title:           req.Title,
		TitleOptions:    req.TitleOptions,
		Summary:         req.Summary,
		MarkdownContent: req.MarkdownContent,
		CoverImageURL:   req.CoverImageURL,
	}
	if err := global.DB.WithContext(ctx).Create(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *SkillArticleService) List(ctx context.Context, req dto.ListSkillArticlesReq) (*vo.SkillArticlePageResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 8
	}
	if req.PageSize > 50 {
		req.PageSize = 50
	}
	query := global.DB.WithContext(ctx).Model(&model.SkillGeneratedArticle{})
	if req.Keyword != "" {
		query = query.Where("keyword = ?", req.Keyword)
	}
	if req.Platform != "" {
		query = query.Where("platform = ?", req.Platform)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []model.SkillGeneratedArticle
	err := query.Order("created_at desc, id desc").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &vo.SkillArticlePageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *SkillArticleService) Get(ctx context.Context, id uint) (*model.SkillGeneratedArticle, error) {
	var article model.SkillGeneratedArticle
	if err := global.DB.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *SkillArticleService) Update(ctx context.Context, id uint, req dto.UpdateSkillArticleReq) (*model.SkillGeneratedArticle, error) {
	article, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Title) != "" {
		article.Title = req.Title
	}
	article.TitleOptions = req.TitleOptions
	article.Summary = req.Summary
	if strings.TrimSpace(req.MarkdownContent) != "" {
		article.MarkdownContent = req.MarkdownContent
	}
	article.CoverImageURL = req.CoverImageURL
	if err := global.DB.WithContext(ctx).Save(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func normalizeTaskID(taskID string) string {
	if strings.TrimSpace(taskID) != "" {
		return strings.TrimSpace(taskID)
	}
	return utils.GenerateUniqueId()
}

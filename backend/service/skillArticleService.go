package service

import (
	"AutoArticle/api/skillArticle/dto"
	"AutoArticle/api/skillArticle/vo"
	"AutoArticle/global"
	"AutoArticle/model"
	"AutoArticle/utils"
	"context"
	"encoding/json"
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

const (
	skillArticleStatusDraft        = "draft"
	skillArticleStatusReady        = "ready_to_publish"
	skillArticlePublishUnpublished = "unpublished"
	skillArticlePublishReady       = "ready"
	skillArticleHumanized          = "done"
)

type SkillArticleService struct{}

type skillPublishPayload struct {
	ArticleID       uint   `json:"articleId"`
	TaskID          string `json:"taskId"`
	Platform        string `json:"platform"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	MarkdownContent string `json:"markdownContent"`
	CoverImageURL   string `json:"coverImageUrl"`
	Images          string `json:"images"`
	Sources         string `json:"sources"`
	Tags            string `json:"tags"`
}

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
		KeywordSegments: req.KeywordSegments,
		Category:        req.Category,
		Title:           req.Title,
		TitleOptions:    req.TitleOptions,
		Summary:         req.Summary,
		MarkdownContent: req.MarkdownContent,
		HTMLPreview:     skillMarkdownToHTML(req.MarkdownContent),
		Tags:            req.Tags,
		CoverImageURL:   req.CoverImageURL,
		CoverImageType:  normalizeDefault(req.CoverImageType, "missing"),
		Images:          req.Images,
		Sources:         req.Sources,
		HotTopics:       req.HotTopics,
		StyleProfile:    req.StyleProfile,
		WordCount:       req.WordCount,
		ModelProvider:   req.ModelProvider,
		ModelName:       req.ModelName,
		PromptVersion:   req.PromptVersion,
		SkillVersion:    req.SkillVersion,
		HumanizeStatus:  normalizeDefault(req.HumanizeStatus, skillArticleHumanized),
		Status:          normalizeDefault(req.Status, skillArticleStatusDraft),
		PublishStatus:   normalizeDefault(req.PublishStatus, skillArticlePublishUnpublished),
		PublishPayload:  req.PublishPayload,
		ErrorMessage:    req.ErrorMessage,
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
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
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
		article.HTMLPreview = skillMarkdownToHTML(req.MarkdownContent)
	}
	article.Tags = req.Tags
	article.CoverImageURL = req.CoverImageURL
	article.CoverImageType = req.CoverImageType
	article.Images = req.Images
	article.Sources = req.Sources
	article.StyleProfile = req.StyleProfile
	article.WordCount = req.WordCount
	article.HumanizeStatus = req.HumanizeStatus
	article.Status = req.Status
	article.PublishStatus = req.PublishStatus
	article.PublishPayload = req.PublishPayload
	article.ErrorMessage = req.ErrorMessage
	if err := global.DB.WithContext(ctx).Save(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func (s *SkillArticleService) CreatePublishPackage(ctx context.Context, id uint) (*vo.SkillPublishPackageResp, error) {
	article, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	payload := skillPublishPayload{
		ArticleID:       article.ID,
		TaskID:          article.TaskID,
		Platform:        article.Platform,
		Title:           article.Title,
		Summary:         article.Summary,
		MarkdownContent: article.MarkdownContent,
		CoverImageURL:   article.CoverImageURL,
		Images:          article.Images,
		Sources:         article.Sources,
		Tags:            article.Tags,
	}
	payloadBytes, _ := json.Marshal(payload)
	err = global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Updates(map[string]interface{}{
			"status":          skillArticleStatusReady,
			"publish_status":  skillArticlePublishReady,
			"publish_payload": string(payloadBytes),
		}).Error
	})
	if err != nil {
		return nil, err
	}
	article.Status = skillArticleStatusReady
	article.PublishStatus = skillArticlePublishReady
	article.PublishPayload = string(payloadBytes)
	return &vo.SkillPublishPackageResp{Article: *article, Payload: string(payloadBytes)}, nil
}

func normalizeTaskID(taskID string) string {
	if strings.TrimSpace(taskID) != "" {
		return strings.TrimSpace(taskID)
	}
	return utils.GenerateUniqueId()
}

func normalizeDefault(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func skillMarkdownToHTML(markdown string) string {
	lines := strings.Split(markdown, "\n")
	blocks := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "### "):
			blocks = append(blocks, "<h3>"+html.EscapeString(strings.TrimPrefix(trimmed, "### "))+"</h3>")
		case strings.HasPrefix(trimmed, "## "):
			blocks = append(blocks, "<h2>"+html.EscapeString(strings.TrimPrefix(trimmed, "## "))+"</h2>")
		case strings.HasPrefix(trimmed, "# "):
			blocks = append(blocks, "<h1>"+html.EscapeString(strings.TrimPrefix(trimmed, "# "))+"</h1>")
		default:
			blocks = append(blocks, "<p>"+html.EscapeString(trimmed)+"</p>")
		}
	}
	return strings.Join(blocks, "\n")
}

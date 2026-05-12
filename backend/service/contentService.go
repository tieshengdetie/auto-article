package service

import (
	"AutoArticle/api/content/dto"
	"AutoArticle/api/content/vo"
	"AutoArticle/common/llm"
	llmtypes "AutoArticle/common/llm/types"
	"AutoArticle/constants"
	"AutoArticle/global"
	"AutoArticle/model"
	"AutoArticle/prompt"
	"AutoArticle/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	searchTypeChannel  = "channel"
	searchTypeGeneral  = "general"
	articleStatusDraft = "draft"
	articleStatusReady = "ready_to_publish"
	defaultQwenModel   = "qwen3.5-plus"
	defaultDoubaoModel = constants.Doubao20Pro260115
	llmProviderQwen    = "qwen"
	llmProviderDoubao  = "doubao"
)

type ContentService struct{}

type articlePromptNews struct {
	Title       string
	Source      string
	Ctime       string
	URL         string
	Description string
}

type articlePromptData struct {
	Keyword    string
	UserPrompt string
	News       []articlePromptNews
}

type generatedArticlePayload struct {
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	MarkdownContent string   `json:"markdownContent"`
	Tags            []string `json:"tags"`
}

type articleLLMRequest struct {
	Provider     string
	ModelName    string
	SystemPrompt string
	UserPrompt   string
}

type articleLLMResponse struct {
	Content   string
	ModelName string
	Provider  string
}

func (s *ContentService) SaveNews(ctx context.Context, req dto.SaveNewsReq) (*vo.SaveNewsResp, error) {
	if strings.TrimSpace(req.Keyword) == "" {
		req.Keyword = "未指定关键词"
	}
	if len(req.List) == 0 {
		return nil, errors.New("新闻列表不能为空")
	}
	if req.SearchType == "" {
		req.SearchType = searchTypeGeneral
	}

	resp := &vo.SaveNewsResp{List: make([]model.CollectedNews, 0, len(req.List))}
	for _, item := range req.List {
		if strings.TrimSpace(item.Title) == "" && strings.TrimSpace(item.URL) == "" {
			resp.Skipped++
			continue
		}
		news := model.CollectedNews{
			UniqueKey:    buildNewsUniqueKey(req.Keyword, req.SearchType, req.ChannelID, item),
			Keyword:      req.Keyword,
			SearchType:   req.SearchType,
			ChannelID:    req.ChannelID,
			ChannelName:  req.ChannelName,
			ThirdPartyID: item.ID,
			URL:          item.URL,
			Ctime:        item.Ctime,
			Title:        item.Title,
			PicURL:       item.PicURL,
			Source:       item.Source,
			Description:  item.Description,
		}
		result := global.DB.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoNothing: true,
		}).Create(&news)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			resp.Skipped++
			var existing model.CollectedNews
			if err := global.DB.WithContext(ctx).Where("unique_key = ?", news.UniqueKey).First(&existing).Error; err == nil {
				resp.List = append(resp.List, existing)
			}
			continue
		}
		resp.Saved++
		resp.List = append(resp.List, news)
	}
	return resp, nil
}

func (s *ContentService) ListNews(ctx context.Context, req dto.ListNewsReq) ([]model.CollectedNews, error) {
	query := global.DB.WithContext(ctx).Model(&model.CollectedNews{})
	if req.Keyword != "" {
		query = query.Where("keyword = ?", req.Keyword)
	}
	if req.SearchType != "" {
		query = query.Where("search_type = ?", req.SearchType)
	}
	if req.ChannelID > 0 {
		query = query.Where("channel_id = ?", req.ChannelID)
	}
	var list []model.CollectedNews
	err := query.Order("id desc").Limit(100).Find(&list).Error
	return list, err
}

func (s *ContentService) GenerateArticle(ctx context.Context, req dto.GenerateArticleReq) (*model.GeneratedArticle, error) {
	if strings.TrimSpace(req.Keyword) == "" {
		return nil, errors.New("关键词不能为空")
	}
	if len(req.NewsIDs) == 0 {
		return nil, errors.New("请选择用于生成文章的新闻")
	}

	var newsList []model.CollectedNews
	if err := global.DB.WithContext(ctx).Where("id IN ?", req.NewsIDs).Find(&newsList).Error; err != nil {
		return nil, err
	}
	if len(newsList) == 0 {
		return nil, errors.New("未找到可用新闻")
	}

	systemPrompt, err := renderArticleSystemPrompt(req.Keyword, req.UserPrompt, newsList)
	if err != nil {
		return nil, err
	}
	llmResp, err := callArticleLLM(ctx, articleLLMRequest{
		Provider:     req.LLMProvider,
		ModelName:    req.ModelName,
		SystemPrompt: systemPrompt,
		UserPrompt:   buildArticleUserPrompt(req.UserPrompt),
	})
	if err != nil {
		return nil, err
	}

	payload, err := parseArticlePayload(llmResp.Content)
	if err != nil {
		return nil, err
	}
	sourceIDs, _ := json.Marshal(req.NewsIDs)
	coverURL, imageStatus := pickCoverImage(newsList)
	article := model.GeneratedArticle{
		Keyword:         req.Keyword,
		Title:           payload.Title,
		Summary:         payload.Summary,
		MarkdownContent: payload.MarkdownContent,
		HTMLPreview:     markdownToHTML(payload.MarkdownContent),
		Tags:            strings.Join(payload.Tags, ","),
		CoverImageURL:   coverURL,
		ImageStatus:     imageStatus,
		ModelName:       llmResp.ModelName,
		SourceNewsIDs:   string(sourceIDs),
		Status:          articleStatusDraft,
	}
	if err := global.DB.WithContext(ctx).Create(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *ContentService) ListArticles(ctx context.Context, req dto.ListArticlesReq) (*vo.ArticlePageResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 8
	}
	if req.PageSize > 50 {
		req.PageSize = 50
	}
	query := global.DB.WithContext(ctx).Model(&model.GeneratedArticle{})
	if req.Keyword != "" {
		query = query.Where("keyword = ?", req.Keyword)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []model.GeneratedArticle
	err := query.Order("created_at desc, id desc").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &vo.ArticlePageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *ContentService) GetArticle(ctx context.Context, id uint) (*model.GeneratedArticle, error) {
	var article model.GeneratedArticle
	if err := global.DB.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *ContentService) UpdateArticle(ctx context.Context, id uint, req dto.UpdateArticleReq) (*model.GeneratedArticle, error) {
	article, err := s.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Title) != "" {
		article.Title = req.Title
	}
	article.Summary = req.Summary
	article.MarkdownContent = req.MarkdownContent
	article.HTMLPreview = markdownToHTML(req.MarkdownContent)
	article.CoverImageURL = req.CoverImageURL
	article.Tags = req.Tags
	if article.CoverImageURL == "" {
		article.ImageStatus = "missing"
	} else {
		article.ImageStatus = "reused"
	}
	if err := global.DB.WithContext(ctx).Save(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func (s *ContentService) CreatePublishPackage(ctx context.Context, articleID uint) (*vo.PublishPackageResp, error) {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}
	pkg := model.PublishPackage{
		ArticleID:       article.ID,
		Platform:        "toutiao",
		Title:           article.Title,
		MarkdownContent: article.MarkdownContent,
		CoverImageURL:   article.CoverImageURL,
		SourceNewsIDs:   article.SourceNewsIDs,
		Status:          articleStatusReady,
		Remark:          "头条号真实发布接口未接入，请复制待发布包内容到头条号后台发布。",
	}
	err = global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pkg).Error; err != nil {
			return err
		}
		return tx.Model(article).Update("status", articleStatusReady).Error
	})
	if err != nil {
		return nil, err
	}
	article.Status = articleStatusReady
	return &vo.PublishPackageResp{Package: pkg, Article: *article}, nil
}

func buildNewsUniqueKey(keyword, searchType string, channelID int, item dto.NewsItemReq) string {
	raw := fmt.Sprintf("%s|%s|%d|%s|%s|%s", keyword, searchType, channelID, item.ID, item.URL, item.Title)
	return utils.Md5(raw)
}

func renderArticleSystemPrompt(keyword, userPrompt string, newsList []model.CollectedNews) (string, error) {
	news := make([]articlePromptNews, 0, len(newsList))
	for _, item := range newsList {
		news = append(news, articlePromptNews{
			Title:       item.Title,
			Source:      item.Source,
			Ctime:       item.Ctime,
			URL:         item.URL,
			Description: item.Description,
		})
	}
	return prompt.Global().Render("article_generate", articlePromptData{
		Keyword:    keyword,
		UserPrompt: strings.TrimSpace(userPrompt),
		News:       news,
	})
}

func callArticleLLM(ctx context.Context, req articleLLMRequest) (*articleLLMResponse, error) {
	provider := strings.ToLower(strings.TrimSpace(req.Provider))
	if provider == "" {
		provider = llmProviderQwen
	}
	switch provider {
	case llmProviderQwen:
		return callArticleWithQwen(ctx, req)
	case llmProviderDoubao:
		return callArticleWithDoubao(ctx, req)
	default:
		return nil, fmt.Errorf("不支持的大模型供应商: %s", req.Provider)
	}
}

func callArticleWithQwen(ctx context.Context, req articleLLMRequest) (*articleLLMResponse, error) {
	if global.ServerConfig.AliYunModel.AppKey == "" {
		return nil, errors.New("未配置千问 API Key")
	}
	modelName := strings.TrimSpace(req.ModelName)
	if modelName == "" {
		modelName = defaultQwenModel
	}
	client := llm.NewQwenClimberClient(global.ServerConfig.AliYunModel.AppKey)
	choices, _, err := client.CallQwen(ctx, llmtypes.QwenOption{
		ModelName:     modelName,
		SystemMessage: []string{req.SystemPrompt},
		UserMessage:   req.UserPrompt,
		Temperature:   0.7,
		MaxTokens:     16384,
	})
	if err != nil {
		return nil, err
	}
	if len(choices) == 0 || strings.TrimSpace(choices[0].Message.Content) == "" {
		return nil, errors.New("千问返回内容为空")
	}
	return &articleLLMResponse{
		Content:   choices[0].Message.Content,
		ModelName: modelName,
		Provider:  llmProviderQwen,
	}, nil
}

func callArticleWithDoubao(ctx context.Context, req articleLLMRequest) (*articleLLMResponse, error) {
	if global.ServerConfig.VolcengineModel.AppKey == "" {
		return nil, errors.New("未配置豆包 API Key")
	}
	modelName := strings.TrimSpace(req.ModelName)
	if modelName == "" {
		modelName = defaultDoubaoModel
	}
	client := llm.NewDBaoClimberClient(global.ServerConfig.VolcengineModel.AppKey, global.ServerConfig.VolcengineModel.BaseUrl)
	resp, _, err := client.CallDBao(ctx, llmtypes.DBaoOption{
		ModelName:           modelName,
		EndPointId:          modelName,
		SystemMessage:       []string{req.SystemPrompt},
		UserMessage:         req.UserPrompt,
		Temperature:         0.7,
		MaxCompletionTokens: 16384,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Choices) == 0 || resp.Choices[0].Message.Content.StringValue == nil {
		return nil, errors.New("豆包返回内容为空")
	}
	content := strings.TrimSpace(*resp.Choices[0].Message.Content.StringValue)
	if content == "" {
		return nil, errors.New("豆包返回内容为空")
	}
	return &articleLLMResponse{
		Content:   content,
		ModelName: modelName,
		Provider:  llmProviderDoubao,
	}, nil
}

func buildArticleUserPrompt(userPrompt string) string {
	if strings.TrimSpace(userPrompt) != "" {
		return "请严格结合系统提示词中的关键词、新闻资料和以下用户要求生成文章：" + userPrompt
	}
	return "请严格结合系统提示词中的关键词和新闻资料，润色整合成一篇新的自媒体文章。"
}

func parseArticlePayload(content string) (*generatedArticlePayload, error) {
	cleaned := strings.TrimSpace(content)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)
	var payload generatedArticlePayload
	if err := json.Unmarshal([]byte(cleaned), &payload); err != nil {
		return nil, fmt.Errorf("解析大模型文章 JSON 失败: %w", err)
	}
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.MarkdownContent) == "" {
		return nil, errors.New("大模型返回缺少标题或正文")
	}
	return &payload, nil
}

func pickCoverImage(newsList []model.CollectedNews) (string, string) {
	for _, item := range newsList {
		if item.PicURL != "" {
			return item.PicURL, "reused"
		}
	}
	return "", "missing"
}

func markdownToHTML(markdown string) string {
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

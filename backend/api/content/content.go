package content

import (
	"AutoArticle/api/content/dto"
	"AutoArticle/service"
	"AutoArticle/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Content struct {
	contentService *service.ContentService
}

func NewContent() *Content {
	return &Content{contentService: &service.ContentService{}}
}

func (c *Content) HandleSaveNews(ctx *gin.Context) {
	var req dto.SaveNewsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := c.contentService.SaveNews(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleListNews(ctx *gin.Context) {
	var req dto.ListNewsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := c.contentService.ListNews(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleGenerateArticle(ctx *gin.Context) {
	var req dto.GenerateArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := c.contentService.GenerateArticle(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleListArticles(ctx *gin.Context) {
	var req dto.ListArticlesReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := c.contentService.ListArticles(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleGetArticle(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	resp, err := c.contentService.GetArticle(ctx.Request.Context(), id)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleUpdateArticle(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var req dto.UpdateArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := c.contentService.UpdateArticle(ctx.Request.Context(), id, req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (c *Content) HandleCreatePublishPackage(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	resp, err := c.contentService.CreatePublishPackage(ctx.Request.Context(), id)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func parseID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Err(ctx, 400, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

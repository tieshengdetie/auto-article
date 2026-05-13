package skillArticle

import (
	"AutoArticle/api/skillArticle/dto"
	"AutoArticle/service"
	"AutoArticle/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SkillArticle struct {
	skillArticleService *service.SkillArticleService
}

func NewSkillArticle() *SkillArticle {
	return &SkillArticle{skillArticleService: &service.SkillArticleService{}}
}

func (s *SkillArticle) HandleCreate(ctx *gin.Context) {
	var req dto.CreateSkillArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := s.skillArticleService.Create(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (s *SkillArticle) HandleList(ctx *gin.Context) {
	var req dto.ListSkillArticlesReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := s.skillArticleService.List(ctx.Request.Context(), req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (s *SkillArticle) HandleGet(ctx *gin.Context) {
	id, ok := parseSkillArticleID(ctx)
	if !ok {
		return
	}
	resp, err := s.skillArticleService.Get(ctx.Request.Context(), id)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (s *SkillArticle) HandleUpdate(ctx *gin.Context) {
	id, ok := parseSkillArticleID(ctx)
	if !ok {
		return
	}
	var req dto.UpdateSkillArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Err(ctx, 400, err.Error())
		return
	}
	resp, err := s.skillArticleService.Update(ctx.Request.Context(), id, req)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func (s *SkillArticle) HandleCreatePublishPackage(ctx *gin.Context) {
	id, ok := parseSkillArticleID(ctx)
	if !ok {
		return
	}
	resp, err := s.skillArticleService.CreatePublishPackage(ctx.Request.Context(), id)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	utils.Success(ctx, resp)
}

func parseSkillArticleID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Err(ctx, 400, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

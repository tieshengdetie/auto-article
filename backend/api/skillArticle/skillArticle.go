package skillArticle

import (
	"AutoArticle/api/skillArticle/dto"
	"AutoArticle/service"
	"AutoArticle/utils"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func (s *SkillArticle) HandleUploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		utils.Err(ctx, 400, "请选择要上传的图片")
		return
	}
	if file.Size <= 0 || file.Size > 8*1024*1024 {
		utils.Err(ctx, 400, "图片大小需在 8MB 以内")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	allowedExt := map[string]bool{".jpg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowedExt[ext] {
		utils.Err(ctx, 400, "仅支持 jpg、png、webp、gif 图片")
		return
	}

	now := time.Now()
	relativeDir := filepath.Join("article-images", "uploads", now.Format("2006"), now.Format("01"))
	targetDir := filepath.Join("static", relativeDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}

	name, err := randomImageName(ext)
	if err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}
	targetPath := filepath.Join(targetDir, name)
	if err := ctx.SaveUploadedFile(file, targetPath); err != nil {
		utils.Err(ctx, 500, err.Error())
		return
	}

	publicPath := "/" + filepath.ToSlash(filepath.Join("static", relativeDir, name))
	utils.Success(ctx, gin.H{
		"url":      publicPath,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

func parseSkillArticleID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Err(ctx, 400, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

func randomImageName(ext string) (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), hex.EncodeToString(bytes[:]), ext), nil
}

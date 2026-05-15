package router

import (
	"AutoArticle/api/skillArticle"

	"github.com/gin-gonic/gin"
)

func InitSkillArticleRouter(Router *gin.RouterGroup) {
	SkillArticleRouter := Router.Group("skill-articles")
	skillArticleInstance := skillArticle.NewSkillArticle()

	{
		SkillArticleRouter.POST("", skillArticleInstance.HandleCreate)
		SkillArticleRouter.POST("/upload-image", skillArticleInstance.HandleUploadImage)
		SkillArticleRouter.GET("", skillArticleInstance.HandleList)
		SkillArticleRouter.GET("/:id", skillArticleInstance.HandleGet)
		SkillArticleRouter.PUT("/:id", skillArticleInstance.HandleUpdate)
	}
}

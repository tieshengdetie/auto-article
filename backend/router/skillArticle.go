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
		SkillArticleRouter.GET("", skillArticleInstance.HandleList)
		SkillArticleRouter.GET("/:id", skillArticleInstance.HandleGet)
		SkillArticleRouter.PUT("/:id", skillArticleInstance.HandleUpdate)
		SkillArticleRouter.POST("/:id/publish-package", skillArticleInstance.HandleCreatePublishPackage)
	}
}

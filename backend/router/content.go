package router

import (
	"AutoArticle/api/content"

	"github.com/gin-gonic/gin"
)

func InitContentRouter(Router *gin.RouterGroup) {
	ContentRouter := Router.Group("content")
	contentInstance := content.NewContent()

	{
		ContentRouter.POST("/news/save", contentInstance.HandleSaveNews)
		ContentRouter.GET("/news", contentInstance.HandleListNews)
		ContentRouter.POST("/articles/generate", contentInstance.HandleGenerateArticle)
		ContentRouter.GET("/articles", contentInstance.HandleListArticles)
		ContentRouter.GET("/articles/:id", contentInstance.HandleGetArticle)
		ContentRouter.PUT("/articles/:id", contentInstance.HandleUpdateArticle)
		ContentRouter.POST("/articles/:id/publish-package", contentInstance.HandleCreatePublishPackage)
	}
}

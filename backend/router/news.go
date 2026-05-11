package router

import (
	"AutoArticle/api/news"
	"github.com/gin-gonic/gin"
)

// InitNewsRouter 初始化新闻路由
func InitNewsRouter(Router *gin.RouterGroup) {
	NewsRouter := Router.Group("news")
	newsInstance := news.NewNews()

	{
		// 指定分类新闻
		NewsRouter.POST("/all", newsInstance.HandleGetAllNews)
		// 综合新闻
		NewsRouter.POST("/general", newsInstance.HandleGetGeneralNews)
		// 获取新闻频道列表
		NewsRouter.GET("/channels", newsInstance.HandleGetNewsChannels)
	}
}

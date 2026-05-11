package router

import (
	"AutoArticle/api/hotSearch"
	"github.com/gin-gonic/gin"
)

// InitHotSearchRouter 初始化热搜路由
func InitHotSearchRouter(Router *gin.RouterGroup) {
	HotSearchRouter := Router.Group("hotSearch")
	hotSearchInstance := hotSearch.NewHotSearch()

	{
		// 热搜接口
		HotSearchRouter.POST("/", hotSearchInstance.HandleGetHotSearch)
		HotSearchRouter.GET("/types", hotSearchInstance.HandleGetHotSearchTypes)
	}
}

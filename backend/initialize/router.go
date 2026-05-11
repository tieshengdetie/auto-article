package initialize

import (
	"AutoArticle/middleware"
	"AutoArticle/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 注册中间件
	Router.Use(
		middleware.CorsMiddleWare(),    // 跨域的
		middleware.LoggerMiddleWare(),  // 日志
		middleware.RecoverMiddleWare(), // 异常的
	)
	// 配置全局路径
	ApiGroup := Router.Group("/api/v1")
	// 注册路由
	router.InitHotSearchRouter(ApiGroup) // 热搜接口
	router.InitNewsRouter(ApiGroup)      // 新闻接口
	return Router
}

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func syncResponse(ctx *gin.Context, httpStatus int, code int, message string, data interface{}) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  message,
		"data": data,
	})
}

// SyncSuccess 成功的请求
func SyncSuccess(ctx *gin.Context, data interface{}) {
	syncResponse(ctx, http.StatusOK, 10000, "请求成功", data)
}

// SyncFail 失败的请求
func SyncFail(ctx *gin.Context, message string) {
	syncResponse(ctx, http.StatusOK, 0, message, nil)
}

// SyncErr 带错误状态码的请求
func SyncErr(ctx *gin.Context, httpStatus int, message string) {
	syncResponse(ctx, httpStatus, 0, message, nil)
}

type SyncPageVo struct {
	Data       interface{} `json:"data"`       // 数据
	Total      int64       `json:"total"`      // 总条数
	PageSize   int64       `json:"pageSize"`   // 当前条数
	PageNumber int64       `json:"pageNumber"` // 当前页数
}

// SyncBuildPageData 构造分页查询器
func SyncBuildPageData(ctx *gin.Context, data interface{}, total int64) {
	size, number := GetQueryPage(ctx.Request)
	SyncSuccess(ctx, SyncPageVo{
		Data:       If(data == nil, make([]interface{}, 0), data),
		Total:      total,
		PageSize:   size,
		PageNumber: number,
	})
}

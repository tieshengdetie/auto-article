package middleware

import (
	"AutoArticle/global"
	"AutoArticle/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RedisRequestLockRequestLock(key string, timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetInt64("accountId")
		if userId <= 0 {
			utils.Fail(ctx, "redis访问锁要先用户登录")
			ctx.Abort()
			return
		}

		lKey := fmt.Sprintf("lock:%s:%d", key, userId)
		lValue := fmt.Sprintf("%d", time.Now().UnixNano())
		l := utils.NewRedisLock(global.RedisDb, lKey, lValue, timeout)
		err := l.TryLock()
		if err != nil {
			utils.Fail(ctx, "访问评率太高")
			ctx.Abort()
			return
		}
		fmt.Println("进来了-----", lKey, lValue)
		ctx.Next()
		_ = l.UnLock()
	}
}

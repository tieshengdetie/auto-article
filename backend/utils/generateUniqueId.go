package utils

import (
	"AutoArticle/global"
	"context"
	"fmt"
	"time"
)

// GenerateUniqueID
//
//	@Description:  生成唯一ID(依赖redis计数器,如果redis服务出问题,没法生成)
//	@param ctx
//	@return string
//	@return error
func GenerateUniqueID(ctx context.Context) (int64, error) {
	// 获取当前时间戳（毫秒）
	//timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	// 获取当前时间戳（秒）
	timestamp := time.Now().Unix()
	// 使用 Redis 计数器
	counterKey := fmt.Sprintf("counter:%d", timestamp)
	counter, err := global.RedisDb.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, err
	}

	// 设置过期时间，防止计数器无限增长
	_, err = global.RedisDb.Expire(ctx, counterKey, 10*time.Second).Result()
	if err != nil {
		return 0, err
	}

	// 构建唯一 ID
	id := timestamp*100000 + counter

	return id, nil
}

// GenerateUniqueStringID
//
//	@Description:  生成字符串类型的id
//	@param ctx
//	@param redisKey
//	@return string
//	@return error
func GenerateUniqueStringID(ctx context.Context, redisKey string) (string, error) {
	// 获取当前时间戳（秒）
	timestamp := time.Now().Unix()
	// 使用 Redis 计数器
	counterKey := fmt.Sprintf("%v:counter:%d", redisKey, timestamp)
	counter, err := global.RedisDb.Incr(ctx, counterKey).Result()
	if err != nil {
		return "", err
	}

	// 设置过期时间，防止计数器无限增长
	_, err = global.RedisDb.Expire(ctx, counterKey, 10*time.Second).Result()
	if err != nil {
		return "", err
	}

	// 使用 %06d 格式化计数器，确保它至少有 6 位，前面用 0 填充
	formattedCounter := fmt.Sprintf("%06d", counter)

	// 组合时间戳和格式化的计数器生成唯一 ID，而不使用分隔符
	uniqueID := fmt.Sprintf("%d%s", timestamp, formattedCounter)

	return uniqueID, nil
}

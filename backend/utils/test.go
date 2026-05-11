package utils

import (
	"AutoArticle/global"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

func init() {
	idGenerator = &IDGenerator{
		client:      global.RedisDb,
		key:         "counter",
		segmentSize: 200,
	}
}

type IDGenerator struct {
	client      *redis.Client
	key         string
	segmentSize int64 // 一批id的偏移量
	currentID   int64 // 当前可使用的id
	maxID       int64 // 这一批可使用的最大id,下一批从maxID+1开始
	updating    int32 // 使用来表示是否正在更新,updating=1表示由协程在更新新一批id、
	lastUpdate  int64 // 上次更新的时间戳
	initialized int32 // 表示是否已初始化:
}

var idGenerator *IDGenerator

// 生成唯一 ID
func (g *IDGenerator) GenerateUniqueID(ctx context.Context) (string, error) {

	// 1.检查并执行初始化 —— 采取懒加载，第一次获取id才初始化id段
	if atomic.LoadInt32(&g.initialized) == 0 {
		if atomic.CompareAndSwapInt32(&g.initialized, 0, 1) {
			if err := g.getNewSegment(ctx); err != nil {
				// 如果初始化失败，将 initialized 重置为未初始化状态
				atomic.StoreInt32(&g.initialized, 0)
				return "", err
			}
		} else {
			// 等待其他 goroutine 初始化完成
			for atomic.LoadInt32(&g.initialized) == 0 {
				time.Sleep(time.Microsecond)
			}
		}
	}

	// 2. 获取id
	for {
		currentID := atomic.LoadInt64(&g.currentID)
		maxID := atomic.LoadInt64(&g.maxID)

		// 没有可用id: 获取新一批id
		if currentID > maxID {
			if err := g.getNewSegment(ctx); err != nil {
				return "", err
			}
			continue
		}
		// 当前id修改成功才算是获取成功
		if atomic.CompareAndSwapInt64(&g.currentID, currentID, currentID+1) {
			// 如果剩余可用id少于 segmentSize*0.1 就开启异步协程加载新一批id
			if currentID > maxID-int64(float64(g.segmentSize)*0.1) {
				// 只是尝试更新id分区 —— 如果这次CAS失败就等下一个请求获取id时再更新
				go g.tryAsyncUpdate(ctx)
			}
			return fmt.Sprintf("%d", currentID), nil
		}
		// 当前id修改失败就等下轮再获取
		time.Sleep(time.Microsecond)
	}
}

// 获取新的 ID 段
func (g *IDGenerator) getNewSegment(ctx context.Context) error {
	// 1.动态调整分段大小 —— 为避免频繁向redis发送请求,如果id消耗太快就增加分区容量
	// 第一次初始化时lastUpdate=0, 不会受影响
	now := time.Now().UnixNano()
	if now-atomic.LoadInt64(&g.lastUpdate) > int64(time.Second) {
		g.segmentSize = g.segmentSize*2 + 100
		fmt.Printf("将id获取分区容量增加值: %d\n", g.segmentSize)
	}
	// 2.将redis中该key自增到分区号，模拟取了这么多个id
	// 下一批 ID 可用的起始值是 newMaxID + 1
	newMaxID, err := g.client.IncrBy(ctx, g.key, g.segmentSize).Result()
	if err != nil {
		return fmt.Errorf("error fetching new segment: %v", err)
	}

	// 3.atomic确保在更新 ID 段时对共享变量的写入是线程安全的，取消sync.Mutex的使用
	// 初始化时才重置 currentID、后续currentID只有被用了才会修改
	if atomic.LoadInt32(&g.initialized) == 0 {
		atomic.StoreInt64(&g.currentID, newMaxID-g.segmentSize+1)
	}
	// maxID: 这一批可使用的最大id
	atomic.StoreInt64(&g.maxID, newMaxID)
	// lastUpdate: 更新的时间戳
	atomic.StoreInt64(&g.lastUpdate, time.Now().UnixNano())
	fmt.Printf("Fetched new ID segment: currentID=%d, maxID=%d\n", g.currentID, g.maxID)
	return nil
}

// 尝试异步获取一批新的id
func (g *IDGenerator) tryAsyncUpdate(ctx context.Context) {

	// 用cas代替锁: updating=1表示有携程正在更新新一批id
	if atomic.CompareAndSwapInt32(&g.updating, 0, 1) {
		defer atomic.StoreInt32(&g.updating, 0)

		if err := g.getNewSegment(ctx); err != nil {
			fmt.Printf("Error in async fetch: %v\n", err)
		}
	}
}

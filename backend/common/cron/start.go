package cron

import "sync"

type schedulerHandler interface {
	HandleSchedulerWork()
}

type registerHandlerStruct struct {
	spec    string
	handler schedulerHandler
}

type handlerMap map[string]registerHandlerStruct

var handlerMapSyncLock sync.Mutex
var handlerMapObj = make(handlerMap)
var GlobalSchedulerMap = make(map[string]schedulerHandler)

func RegisterFunc() Scheduler {
	handlerMapSyncLock.Lock()
	defer handlerMapSyncLock.Unlock()
	// 初始化定时任务调度器
	schedulerObj := NewScheduler()
	for _, handlerObj := range handlerMapObj {
		schedulerObj.AddFunc(handlerObj.spec, handlerObj.handler.HandleSchedulerWork)
	}
	// 注册定时任务1，每分钟执行一次
	//scheduler.AddFunc("0 * * * * *", HandleSchedulerWork)
	// 启动定时任务调度器
	schedulerObj.Start()
	return schedulerObj
}

func register(uniqueKey, spec string, handler schedulerHandler) {
	handlerMapSyncLock.Lock()
	defer handlerMapSyncLock.Unlock()
	handlerMapObj[uniqueKey] = registerHandlerStruct{
		spec:    spec,
		handler: handler,
	}
	GlobalSchedulerMap[uniqueKey] = handler
}

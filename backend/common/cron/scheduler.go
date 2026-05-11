package cron

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

// Scheduler 是定时任务调度器的接口
type Scheduler interface {
	AddFunc(spec string, cmd func())
	Start()
	Stop()
}

// scheduler 是定时任务调度器的实现
type scheduler struct {
	s *gocron.Scheduler
}

// NewScheduler 创建一个新的定时任务调度器
func NewScheduler() Scheduler {
	s := gocron.NewScheduler(time.UTC)
	return &scheduler{s: s}
}

// AddFunc 添加一个定时任务
func (s *scheduler) AddFunc(spec string, cmd func()) {
	_, err := s.s.CronWithSeconds(spec).Do(cmd)
	if err != nil {
		log.Fatalf("添加定时任务失败: %v", err)
	}
}

// Start 启动定时任务调度器
func (s *scheduler) Start() {
	s.s.StartAsync()
}

// Stop 停止定时任务调度器
func (s *scheduler) Stop() {
	s.s.Stop()
}

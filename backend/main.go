package main

import (
	_ "AutoArticle/common"
	"AutoArticle/common/cron"
	"AutoArticle/global"
	"AutoArticle/initialize"
	"AutoArticle/prompt"
	"AutoArticle/utils"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	envString string
	cronMode  bool
	taskName  string
)

func init() {
	flag.StringVar(&envString, "envString", "dev", "请输入当前环境配置文件")
	flag.BoolVar(&cronMode, "cronMode", false, "是否是执行定时任务模式")
	flag.StringVar(&taskName, "taskName", "taskName", "定时任务名称")
}

func main() {
	flag.Parse()
	fmt.Println(envString, "当前环境")
	// 初始化配置
	initialize.InitConfig(envString)
	// 初始化数据库
	initialize.InitDB()
	initialize.InitRedisDb()
	//初始化路由
	router := initialize.Routers()
	// 初始化提示词加载器
	err := prompt.InitGlobal("prompt/templates")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 启动定时任务
	if cronMode == true {

		handler, ok := cron.GlobalSchedulerMap[taskName]
		if !ok {
			fmt.Println("定时任务不存在")
			return
		}
		handler.HandleSchedulerWork()
		return
	}
	// 初始化定时任务调度器
	//scheduler := cron.RegisterFunc()
	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	// 获取服务地址
	address, _ := utils.GetServiceIp()
	fmt.Println("本服务的地址为：", address)
	// 启动服务
	fmt.Println(fmt.Sprintf("HTTP服务启动:localhost:%s", PORT))
	// 优雅退出程序
	go func() {
		// 启动服务
		if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
			fmt.Println(fmt.Sprintf("HTTP服务启动失败:%s", err.Error()))
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	// 停止定时任务调度器
	//scheduler.Stop()
}

package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
)

var (
	// 接收端口号
	port string

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "启动http服务,使用方法: app server --port=?",
		Run: func(cmd *cobra.Command, args []string) {
			if port == "" {
				log.Fatalf("port参数不能为空")
			}
			engine := gin.Default()
			if err := engine.Run(":" + port); err != nil {
				log.Fatalf("服务器启动失败, err: %s", err.Error())
			}
		},
	}
)

func init() {
	// 将server命令添加为rootCmd的子命令
	rootCmd.AddCommand(serverCmd)
	// server子命令接收port选项参数
	serverCmd.Flags().StringVar(&port, "port", "", "端口号")
}

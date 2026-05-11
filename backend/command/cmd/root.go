package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dream",
	Short: "dream 为命令行处理数据的工具",
	Long:  `dream 为命令行处理数据的工具,非常的好用！`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello,dream...")
	},
}

func init() {
	rootCmd.PersistentFlags().String("version", "", "版本")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}

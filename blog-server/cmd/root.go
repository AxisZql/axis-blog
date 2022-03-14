package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// 利用cobra配置启动命令
var rootCmd = &cobra.Command{
	Use:   "AXIS-BLOG",
	Short: "AXIS-BLOG COMMAND",
	Long:  "AXIS-BLOG COMMAND",
	Run: func(cmd *cobra.Command, args []string) {
		// 配置默认启动
	},
}

//记录配置文件
var cfgFile string

func init() {
	// 定义全局持久化标志
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
}

func Execute() {
	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("服务启动失败:%v\n", err)
		os.Exit(1)
	}
}

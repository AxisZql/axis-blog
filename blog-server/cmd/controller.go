package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api 服务", //短描述
	Long:  "api 服务", //长描述
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func AxAPi() {

	router := gin.Default()
	Routers(router)

}

// Routers 路由设置
func Routers(r *gin.Engine) {

}

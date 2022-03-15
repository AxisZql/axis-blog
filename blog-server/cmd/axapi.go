package cmd

import (
	"blog-server/common"
	util "blog-server/common/tools"
	ctrl "blog-server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api 服务", //短描述
	Long:  "api 服务", //长描述
	Run: func(cmd *cobra.Command, args []string) {
		AxAPi()
	},
}

func AxAPi() {
	// 运行环境初始化
	common.EnvInit()
	router := gin.Default()
	Routers(router)
	port := fmt.Sprintf(":%d", common.Conf.App.Port)
	server := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 32 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("服务在%s端口启动成功", port))

}

// Routers 路由设置
func Routers(r *gin.Engine) {
	var userAuth ctrl.UserAuth
	admin := r.Group("/admin")
	users := r.Group("/users")

	r.POST("/register", userAuth.Register) //用户注册

	users.Use(util.Auth())
	{
		users.GET("/code", userAuth.SendEmailCode)      //发送邮箱验证码
		users.PUT("/password", userAuth.UpdatePassword) //修改密码
		users.POST("/oauth/weibo", userAuth.WeiboLogin) //微博登陆
		users.POST("/users/oauth/qq", userAuth.QQLogin) //QQ登陆
	}

	admin.Use(util.Auth())
	{
		admin.GET("/users/area", userAuth.ListUserAreas)           //获取用户区域分布
		admin.GET("/users", userAuth.ListUsers)                    //查询用户后台列表
		admin.PUT("/users/password", userAuth.UpdateAdminPassword) //修改管理员密码
	}

}

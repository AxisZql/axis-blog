package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	ctrl.UserAuthHandle
}

func (user *UserAuth) SendEmailCode(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"a": "fuck you"})
}

func (user *UserAuth) ListUserAreas(*gin.Context)       {}
func (user *UserAuth) ListUsers(*gin.Context)           {}
func (user *UserAuth) Register(*gin.Context)            {}
func (user *UserAuth) UpdatePassword(*gin.Context)      {}
func (user *UserAuth) UpdateAdminPassword(*gin.Context) {}
func (user *UserAuth) WeiboLogin(*gin.Context)          {}
func (user *UserAuth) QQLogin(*gin.Context)             {}

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

package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	ctrl.UserInfoHandle
}

func (u *UserInfo) UpdateUserInfo(*gin.Context)    {}
func (u *UserInfo) UpdateUserAvatar(*gin.Context)  {}
func (u *UserInfo) SaveUserEmail(*gin.Context)     {}
func (u *UserInfo) UpdateUserRole(*gin.Context)    {}
func (u *UserInfo) UpdateUserDisable(*gin.Context) {}
func (u *UserInfo) ListOnlineUsers(*gin.Context)   {}
func (u *UserInfo) RemoveOnlineUser(*gin.Context)  {}

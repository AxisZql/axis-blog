package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type FriendLink struct {
	ctrl.FriendLink
}

func (f *FriendLink) ListFriendLinks(*gin.Context)        {}
func (f *FriendLink) ListFriendLinksBack(*gin.Context)    {}
func (f *FriendLink) SaveOrUpdateFriendLink(*gin.Context) {}
func (f *FriendLink) DeleteFriendLink(*gin.Context)       {}

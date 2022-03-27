package service

import (
	"github.com/gin-gonic/gin"
)

type FriendLink struct{}

func (f *FriendLink) ListFriendLinks(*gin.Context)        {}
func (f *FriendLink) ListFriendLinksBack(*gin.Context)    {}
func (f *FriendLink) SaveOrUpdateFriendLink(*gin.Context) {}
func (f *FriendLink) DeleteFriendLink(*gin.Context)       {}

package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 5:15 PM
* @desc: 用户信息模块接口
 */

type UserInfoHandle interface {
	UpdateUserInfo(*gin.Context)    //更新用户信息
	UpdateUserAvatar(*gin.Context)  //更新用户头像
	SaveUserEmail(*gin.Context)     //绑定用户邮箱
	UpdateUserRole(*gin.Context)    //更新用户角色信息
	UpdateUserDisable(*gin.Context) //更新用户禁用信息
	ListOnlineUsers(*gin.Context)   //查看在线用户
	RemoveOnlineUser(*gin.Context)  //下线用户
}

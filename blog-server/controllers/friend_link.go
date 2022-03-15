package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 9:00 PM
* @desc: 友链模块接口
 */

type FriendLink interface {
	ListFriendLinks(*gin.Context)        //查看友链列表
	ListFriendLinksBack(*gin.Context)    //查看后台友链列表
	SaveOrUpdateFriendLink(*gin.Context) //保存或者修改友链
	DeleteFriendLink(*gin.Context)       //删除友链

}

package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:39 PM
* @desc: 用户留言模块
 */

type MessageHandle interface {
	SaveMessage(*gin.Context)         //添加留言
	ListMessage(*gin.Context)         //查看留言列表
	ListMessageBack(*gin.Context)     //查看后台留言列表
	UpdateMessageReview(*gin.Context) //审核留言
	DeleteMessage(*gin.Context)       //删除留言
}

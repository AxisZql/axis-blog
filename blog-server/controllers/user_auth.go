package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 3:12 PM
* @desc: 用户账号模块接口
 */

type UserAuthHandle interface {
	SendEmailCode(*gin.Context)       //发送邮箱验证码
	ListUserAreas(*gin.Context)       //获取用户地区分布
	ListUsers(*gin.Context)           //查询后台用户列表
	Register(*gin.Context)            //用户注册
	UpdatePassword(*gin.Context)      //修改密码
	UpdateAdminPassword(*gin.Context) //修改管理员密码
	WeiboLogin(*gin.Context)          //微博登陆
	QQLogin(*gin.Context)             //QQ登陆
}

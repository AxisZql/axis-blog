package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-16 11:23 PM
* @desc: 登陆模块接口
 */

type LoginHandle interface {
	Login(*gin.Context)
	LoginOut(*gin.Context)
}

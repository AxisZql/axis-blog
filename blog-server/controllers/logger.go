package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:55 PM
* @desc: 日志模块接口
 */

type LoggerHandler interface {
	ListOperationLogs(*gin.Context)   //查看操作日志
	DeleteOperationLogs(*gin.Context) //删除操作日志
}

package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:33 PM
* @desc: 页面管理模块接口
 */

type PageHandle interface {
	DeletePage(*gin.Context)       //删除页面
	SaveOrUpdatePage(*gin.Context) //保存或更新页面
	ListPages(*gin.Context)        //获取页面列表
}

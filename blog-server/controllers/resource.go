package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 7:56 PM
* @desc: 资源模块
 */

type ResourceHandle interface {
	ListResources(*gin.Context)        //查看资源列表
	DeleteResource(*gin.Context)       //删除资源
	SaveOrUpdateResource(*gin.Context) //新增或修改资源
	ListResourceOption(*gin.Context)   //查看角色资源选项
}

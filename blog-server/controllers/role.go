package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 7:49 PM
* @desc: 角色模块
 */

type RoleHandle interface {
	ListUserRoles(*gin.Context)    //查询用户角色选项
	ListRoles(*gin.Context)        //查询角色类别
	SaveOrUpdateRole(*gin.Context) //保存或更新角色
	DeleteRoles(*gin.Context)      //删除角色
}

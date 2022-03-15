package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:47 PM
* @desc: 菜单控制器模块接口
 */

type MenuHandle interface {
	ListMenus(*gin.Context)        //查看菜单lieb
	SaveOrUpdateMenu(*gin.Context) //新增或者修改菜单
	DeleteMenu(*gin.Context)       //删除菜单
	ListMenuOptions(*gin.Context)  //查看菜单选项
	ListUserMenus(*gin.Context)    //查看当前用户菜单
}

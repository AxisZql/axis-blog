package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 7:38 PM
* @desc: 标签模块
 */

type TagHandle interface {
	ListTags(*gin.Context)        //查询标签列表
	ListTagBack(*gin.Context)     //后台查询标签列表
	ListTagBySearch(*gin.Context) //搜索文章标签
	SaveOrUpdateTag(*gin.Context) //添加或者修改标签
	DeleteTag(*gin.Context)       //删除标签
}

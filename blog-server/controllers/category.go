package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 9:20 PM
* @desc: 分类控制器接口
 */

type CategoryHandle interface {
	ListCategories(*gin.Context)         //查看分类列表
	ListCategoriesBack(*gin.Context)     //查看后台分类列表
	ListCategoriesBySearch(*gin.Context) //搜索文章分类
	SaveOrUpdateCategory(*gin.Context)   //添加或修改分类
	DeleteCategories(*gin.Context)       //删除分类
}

package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 9:46 PM
* @desc: 我在控制器接口
 */

type ArticleHandle interface {
	ListArchives(*gin.Context)           //查看文章归档列表
	ListArticles(*gin.Context)           //查看首页文章
	ListArticleBack(*gin.Context)        //查看后台文章
	SaveOrUpdateArticle(*gin.Context)    //添加或修改文章
	UpdateArticleTop(*gin.Context)       //修改文章置顶
	UpdateArticleDelete(*gin.Context)    //恢复或删除文章
	SaveArticleImages(*gin.Context)      //上传文章图片
	DeleteArticle(*gin.Context)          //物理删除文章
	GetArticleBackById(*gin.Context)     //根据id查看后台文章
	GetArticleById(*gin.Context)         //根据id查看文章
	ListArticleByCondition(*gin.Context) //根据条件查询文章
	ListArticleBySearch(*gin.Context)    //搜索文章
	SaveArticleLike(*gin.Context)        //点赞文章
}

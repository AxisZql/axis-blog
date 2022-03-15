package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 9:32 PM
* @desc: 博客信息模块接口
 */

type BlogInfo interface {
	GetBlogHomeInfo(*gin.Context)     //查看博客信息
	GetBlogBackInfo(*gin.Context)     //查看后台信息
	SavePhotoAlbumCover(*gin.Context) //上传博客配置图片
	UpdateWebsiteConfig(*gin.Context) //更新网站配置
	GetWebSiteConfig(*gin.Context)    //获取网站配置
	GetAbout(*gin.Context)            //查看关于我信息
	UpdateAbout(*gin.Context)         //修改关于我信息
	SendVoice(*gin.Context)           //上传语音信息
	Report(*gin.Context)              //上传访客信息
}

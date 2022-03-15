package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:20 PM
* @desc: 相册控制模块接口
 */

type PhotoAlbumHandle interface {
	SavePhotoAlbumCover(*gin.Context)    //上传相册封面
	SaveOrUpdatePhotoAlbum(*gin.Context) //保存或者更新相册
	ListPhotoAlbumBack(*gin.Context)     //查看后台相册列表
	ListPhotoAlbumBackInfo(*gin.Context) //获取后台相册列表信息
	GetPhotoAlbumBackById(*gin.Context)  //根据id获取后台相册信息
	DeletePhotoAlbumById(*gin.Context)   //根据id删除相册
	ListPhotoAlbum(*gin.Context)         //获取相册列表
}

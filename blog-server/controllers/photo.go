package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 8:05 PM
* @desc:照片模块接口
 */

type PhotoHandle interface {
	ListPhotos(*gin.Context)         //获取后台照片列表
	UpdatePhoto(*gin.Context)        //更新照片信息
	SavePhoto(*gin.Context)          //保存照片
	UpdatePhotoAlbum(*gin.Context)   //移动照片相册
	UpdatePhotoDelete(*gin.Context)  //更新照片删除状态
	DeletePhotos(*gin.Context)       //删除照片
	ListPhotoByAlbumId(*gin.Context) //根据相册id查看照片列表
}

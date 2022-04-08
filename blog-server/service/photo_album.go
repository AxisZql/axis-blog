package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"time"
)

type PhotoAlbum struct{}

type reqSavePhotoAlbumCover struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (p *PhotoAlbum) SavePhotoAlbumCover(ctx *gin.Context) {
	var form reqSavePhotoAlbumCover
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	f, _ := form.File.Open()
	extendName := strings.Split(form.File.Filename, ".")
	if len(extendName) != 2 && extendName[1] != "png" && extendName[1] != "gif" && extendName[1] != "jpg" {
		Response(ctx, errorcode.ValidError, nil, false, "不支持的图片格式;仅支持png|gif|jpg格式")
		return
	}
	defer f.Close()
	fileData, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		logger.Error(err2.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	fileMD5 := fmt.Sprintf("%x", md5.Sum(fileData))
	fileName := fileMD5 + "." + extendName[1]
	filePath := common.Conf.App.PhotoDir + fileName
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	imgUrl := fmt.Sprintf("%s:%d/fphotos/%s", common.Conf.App.HostName, common.Conf.App.Port, fileName)
	Response(ctx, errorcode.Fail, imgUrl, true, "操作成功")

}

type reqSaveOrUpdatePhotoAlbum struct {
	ID         int64  `json:"id"`
	AlbumCover string `json:"albumCover"`
	AlbumDesc  string `json:"albumDesc"`
	AlbumLabel string `json:"albumLabel"`
	AlbumName  string `json:"albumName"`
	Status     int    `json:"status"`
}

func (p *PhotoAlbum) SaveOrUpdatePhotoAlbum(ctx *gin.Context) {
	var form reqSaveOrUpdatePhotoAlbum
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.ID != 0 {
		var pa common.TPhotoAlbum
		r1 := db.Where("id = ?", form.ID).First(&pa)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		pa.AlbumCover = form.AlbumCover
		pa.AlbumDesc = form.AlbumDesc
		pa.AlbumName = form.AlbumName
		pa.Status = form.Status
		pa.UpdateTime = time.Now()
		r1 = db.Save(&pa)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		pa := common.TPhotoAlbum{
			AlbumCover: form.AlbumCover,
			AlbumDesc:  form.AlbumDesc,
			AlbumName:  form.AlbumName,
			Status:     form.Status,
		}
		r1 := db.Model(&common.TPhotoAlbum{}).Create(&pa)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")

}

type reqListPhotoAlbumBack struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (p *PhotoAlbum) ListPhotoAlbumBack(ctx *gin.Context) {
	var form reqListPhotoAlbumBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Size <= 0 || form.Current <= 0 {
		form.Size = 10
		form.Current = 1
	}
	db := common.GetGorm()
	var albumList []struct {
		AlbumCover string `json:"albumCover"`
		AlbumDesc  string `json:"albumDesc"`
		AlbumName  string `json:"albumName"`
		ID         int64  `json:"id"`
		Status     int    `json:"status"`
		PhotoCount int64  `json:"photoCount"`
	}
	var count int64
	r1 := db.Table("t_photo_album").Count(&count)
	r1 = db.Table("t_photo_album").Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&albumList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, val := range albumList {
		var _count int64
		r2 := db.Model(&common.TPhoto{}).Where("album_id = ?", val.ID).Count(&_count)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		albumList[i].PhotoCount = _count
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = albumList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

func (p *PhotoAlbum) ListPhotoAlbumBackInfo(ctx *gin.Context) {
	db := common.GetGorm()
	var albumList []struct {
		ID         int64  `json:"id"`
		AlbumCover string `json:"albumCover"`
		AlbumDesc  string `json:"albumDesc"`
		AlbumName  string `json:"albumName"`
	}
	r1 := db.Table("t_photo_album").Find(&albumList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, albumList, true, "操作成功")
}

type reqGetPhotoAlbumBackById struct {
	AlbumsId int64  `uri:"albumsId" binding:"required"`
	Info     string `uri:"info" binding:"required"`
}

func (p *PhotoAlbum) GetPhotoAlbumBackById(ctx *gin.Context) {
	var form reqGetPhotoAlbumBackById
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Info != "info" {
		Response(ctx, errorcode.NotFoundResource, nil, false, "找不到资源")
		return
	}

	db := common.GetGorm()
	var album struct {
		ID         int64  `json:"id"`
		AlbumCover string `json:"albumCover"`
		AlbumDesc  string `json:"albumDesc"`
		AlbumName  string `json:"albumName"`
		Status     int    `json:"status"`
		PhotoCount int64  `json:"photoCount"`
	}
	var count int64
	r1 := db.Table("t_photo_album").Where("id = ?", form.AlbumsId).First(&album)
	r1 = db.Model(&common.TPhoto{}).Where("album_id = ?", form.AlbumsId).Count(&count)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	album.PhotoCount = count
	Response(ctx, errorcode.Success, album, false, "操作成功")
}
func (p *PhotoAlbum) DeletePhotoAlbumById(*gin.Context) {}
func (p *PhotoAlbum) ListPhotoAlbum(*gin.Context)       {}

package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"strconv"
	"time"
)

type Photo struct{}

type reqListPhotos struct {
	Current  int   `form:"current"`
	Size     int   `form:"size"`
	AlbumId  int64 `form:"albumId"`
	IsDelete int   `form:"isDelete"`
}

func (p *Photo) ListPhotos(ctx *gin.Context) {
	var form reqListPhotos
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	var photoList []struct {
		ID        int64  `json:"id"`
		PhotoDesc string `json:"photoDesc"`
		PhotoName string `json:"photoName"`
		PhotoSrc  string `json:"photoSrc"`
	}
	var count int64
	if form.AlbumId != 0 {
		r1 := db.Model(&common.TPhoto{}).Where("is_delete = ? AND album_id = ?", form.IsDelete, form.AlbumId).Count(&count)
		r1 = db.Table("t_photo").Limit(form.Size).Where("is_delete = ? AND album_id = ?", form.IsDelete, form.AlbumId).Offset((form.Current - 1) * form.Size).Find(&photoList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Model(&common.TPhoto{}).Where("is_delete = ?", form.IsDelete).Count(&count)
		r1 = db.Table("t_photo").Limit(form.Size).Where("is_delete = ?", form.IsDelete).Offset((form.Current - 1) * form.Size).Find(&photoList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}

	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = photoList
	Response(ctx, errorcode.Success, data, true, "操作成功")

}

type reqUpdatePhoto struct {
	ID        int64  `json:"id" binding:"required"`
	PhotoDesc string `json:"photoDesc"`
	PhotoName string `json:"photoName"`
	PhotoSrc  string `json:"photoSrc"`
}

func (p *Photo) UpdatePhoto(ctx *gin.Context) {
	var form reqUpdatePhoto
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var ph common.TPhoto
	r1 := db.Model(&common.TPhoto{}).Where("id = ?", form.ID).First(&ph)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	ph.PhotoName = form.PhotoName
	ph.PhotoDesc = form.PhotoDesc
	ph.PhotoSrc = form.PhotoSrc
	ph.UpdateTime = time.Now()
	r2 := db.Save(&ph)
	if r2.Error != nil {
		logger.Error(r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqSavePhoto struct {
	AlbumId      string   `json:"albumId" binding:"required"`
	PhotoUrlList []string `json:"photoUrlList" binding:"required"`
}

func (p *Photo) SavePhoto(ctx *gin.Context) {
	var form reqSavePhoto
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	aId, _ := strconv.Atoi(form.AlbumId)
	if aId <= 0 {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	for _, val := range form.PhotoUrlList {
		ph := common.TPhoto{
			AlbumId:  int64(aId),
			PhotoSrc: val,
		}
		r1 := db.Model(&common.TPhoto{}).Create(&ph)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqUpdatePhotoAlbum struct {
	AlbumId     int64   `json:"albumId" binding:"required"`
	PhotoIdList []int64 `json:"photoIdList"`
}

func (p *Photo) UpdatePhotoAlbum(ctx *gin.Context) {
	var form reqUpdatePhotoAlbum
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range form.PhotoIdList {
		var ph common.TPhoto
		r1 := db.Model(&common.TPhoto{}).Where("id = ?", val).First(&ph)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		ph.AlbumId = form.AlbumId
		ph.UpdateTime = time.Now()
		r2 := db.Save(&ph)
		if r2.Error != nil {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqUpdatePhotoDelete struct {
	IdList   []int64 `json:"idList"`
	IsDelete int     `json:"isDelete"`
}

func (p *Photo) UpdatePhotoDelete(ctx *gin.Context) {
	var form reqUpdatePhotoDelete
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range form.IdList {
		var ph common.TPhoto
		r1 := db.Model(&common.TPhoto{}).Where("id = ?", val).First(&ph)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		ph.IsDelete = form.IsDelete
		ph.UpdateTime = time.Now()
		r2 := db.Save(&ph)
		if r2.Error != nil {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (p *Photo) DeletePhotos(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	str := fmt.Sprintf("%v", string(data))
	var idList []int64
	err := json.Unmarshal([]byte(str), &idList)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	db := common.GetGorm()
	for _, val := range idList {
		r1 := db.Where("id = ?", val).Delete(&common.TPhoto{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (p *Photo) ListPhotoByAlbumId(*gin.Context) {}

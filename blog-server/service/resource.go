package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Resource struct{}

type reqListResources struct {
	Keywords string `form:"keywords"`
}

type resourceListResources struct {
	Children      []resourceListResources `json:"children"`
	CreateTime    time.Time               `json:"createTime"`
	ID            int64                   `json:"id"`
	IsAnonymous   int                     `json:"isAnonymous"`
	IsDisable     int                     `json:"isDisable"`
	RequestMethod string                  `json:"requestMethod"`
	ResourceName  string                  `json:"resourceName"`
	Url           string                  `json:"url"`
}

func (r *Resource) ListResources(ctx *gin.Context) {
	var form reqListResources
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var rList []resourceListResources

	r1 := db.Model(&common.TResource{}).Where(fmt.Sprintf("isNull(parent_id) AND resource_name LIKE %q", "%"+form.Keywords+"%")).Find(&rList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, m := range rList {
		var child []resourceListResources
		r2 := db.Model(&common.TResource{}).Where("parent_id = ?", m.ID).Find(&child)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		rList[i].Children = child
	}
	Response(ctx, errorcode.Success, rList, true, "操作成功")

}
func (r *Resource) DeleteResource(ctx *gin.Context) {
	Response(ctx, errorcode.Success, nil, true, "暂时不开放模块删除功能")
}
func (r *Resource) SaveOrUpdateResource(ctx *gin.Context) {
	Response(ctx, errorcode.Success, nil, true, "暂时不开放模块修改和创建功能")
}

type resourceOptions struct {
	ID           int64             `json:"id"`
	ResourceName string            `json:"label"`
	Children     []resourceOptions `json:"children"`
}

func (r *Resource) ListResourceOption(ctx *gin.Context) {
	db := common.GetGorm()
	var rList []resourceOptions

	r1 := db.Model(&common.TResource{}).Where("isNull(parent_id)").Find(&rList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, m := range rList {
		var child []resourceOptions
		r2 := db.Model(&common.TResource{}).Where("parent_id = ?", m.ID).Find(&child)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		rList[i].Children = child
	}
	Response(ctx, errorcode.Success, rList, true, "操作成功")
}

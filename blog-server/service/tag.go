package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

type Tag struct{}

func (t *Tag) ListTags(*gin.Context) {}

type reqListTagBack struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (t *Tag) ListTagBack(ctx *gin.Context) {
	var form reqListTagBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	var tagList []struct {
		ArticleCount int64     `json:"articleCount"`
		CreateTime   time.Time `json:"createTime"`
		ID           int64     `json:"id"`
		TagName      string    `json:"tagName"`
	}
	var count int64
	db := common.GetGorm()
	if form.Keywords != "" {
		r1 := db.Table("t_tag").Where(fmt.Sprintf("tag_name LIKE %q", "%"+form.Keywords+"%")).Count(&count)
		r1 = db.Table("t_tag").Where(fmt.Sprintf("tag_name LIKE %q", "%"+form.Keywords+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&tagList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Table("t_tag").Count(&count)
		r1 = db.Table("t_tag").Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&tagList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}

	for i, val := range tagList {
		var _count int64
		r2 := db.Model(&common.TArticleTag{}).Where("tag_id = ?", val.ID).Count(&_count)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}

		tagList[i].ArticleCount = _count
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = tagList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqListTagBySearch struct {
	Keywords string `form:"keywords"`
}

func (t *Tag) ListTagBySearch(ctx *gin.Context) {
	var form reqListTagBySearch
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var tagList []struct {
		ID      int64  `json:"id"`
		TagName string `json:"tagName"`
	}
	if form.Keywords == "" {
		r1 := db.Model(&common.TTag{}).Limit(10).Order("create_time DESC").Find(&tagList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		Response(ctx, errorcode.Success, tagList, true, "操作成功")
		return
	} else {
		r1 := db.Model(&common.TTag{}).Where(fmt.Sprintf("tag_name LIKE %q", "%"+form.Keywords+"%")).Limit(10).Order("create_time DESC").Find(&tagList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		Response(ctx, errorcode.Success, tagList, true, "操作成功")
		return
	}
}

type reqSaveOrUpdateTag struct {
	ArticleCount int64  `json:"articleCount"`
	ID           int64  `json:"id"`
	TagName      string `json:"tagName"`
}

func (t *Tag) SaveOrUpdateTag(ctx *gin.Context) {
	var form reqSaveOrUpdateTag
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.ID != 0 {
		var tag common.TTag
		r1 := db.Where("id = ?", form.ID).First(&tag)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		tag.TagName = form.TagName
		tag.UpdateTime = time.Now()
		r1 = db.Save(&tag)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		tag := common.TTag{
			TagName: form.TagName,
		}
		r1 := db.Model(&common.TTag{}).Create(&tag)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")

}
func (t *Tag) DeleteTag(ctx *gin.Context) {
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
		var t common.TArticleTag
		r1 := db.Where("tag_id = ?", val).First(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if r1.Error == nil {
			Response(ctx, errorcode.Fail, nil, false, "该标签正在使用,无法删除")
			return
		} else {
			r2 := db.Where("id = ?", val).Delete(&common.TTag{})
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

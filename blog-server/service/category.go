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

type Category struct{}

func (c *Category) ListCategories(*gin.Context) {}

type reqListCategoriesBack struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (c *Category) ListCategoriesBack(ctx *gin.Context) {
	var form reqListCategoriesBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	var categoryList []struct {
		ArticleCount int64     `json:"articleCount"`
		CreateTime   time.Time `json:"createTime"`
		ID           int64     `json:"id"`
		CategoryName string    `json:"categoryName"`
	}
	var count int64
	db := common.GetGorm()
	if form.Keywords != "" {
		r1 := db.Table("t_category").Where(fmt.Sprintf("category_name LIKE %q", "%"+form.Keywords+"%")).Count(&count)
		r1 = db.Table("t_category").Where(fmt.Sprintf("category_name LIKE %q", "%"+form.Keywords+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&categoryList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Table("t_category").Count(&count)
		r1 = db.Table("t_category").Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&categoryList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}

	for i, val := range categoryList {
		var _count int64
		r2 := db.Model(&common.TArticle{}).Where("category_id = ?", val.ID).Count(&_count)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}

		categoryList[i].ArticleCount = _count
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = categoryList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqListCategoriesBySearch struct {
	Keywords string `form:"keywords"`
}

func (c *Category) ListCategoriesBySearch(ctx *gin.Context) {
	var form reqListCategoriesBySearch
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var categoryList []struct {
		ID           int64  `json:"id"`
		CategoryName string `json:"categoryName"`
	}
	if form.Keywords == "" {
		r1 := db.Model(&common.TCategory{}).Limit(10).Order("create_time DESC").Find(&categoryList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		Response(ctx, errorcode.Success, categoryList, true, "操作成功")
		return
	} else {
		r1 := db.Model(&common.TCategory{}).Where(fmt.Sprintf("category_name LIKE %q", "%"+form.Keywords+"%")).Limit(10).Order("create_time DESC").Find(&categoryList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		Response(ctx, errorcode.Success, categoryList, true, "操作成功")
		return
	}

}

type reqSaveOrUpdateCategory struct {
	CategoryName string `json:"categoryName"`
	ID           int64  `json:"id"`
}

func (c *Category) SaveOrUpdateCategory(ctx *gin.Context) {
	var form reqSaveOrUpdateCategory
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.ID != 0 {
		var category common.TCategory
		r1 := db.Where("id = ?", form.ID).First(&category)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		category.CategoryName = form.CategoryName
		category.UpdateTime = time.Now()
		r1 = db.Save(&category)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		ca := common.TCategory{
			CategoryName: form.CategoryName,
		}
		r1 := db.Model(&common.TCategory{}).Create(&ca)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (c *Category) DeleteCategories(ctx *gin.Context) {
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
		var t common.TArticle
		r1 := db.Where("category_id = ?", val).First(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if r1.Error == nil {
			Response(ctx, errorcode.Fail, nil, false, "该分类正在使用,无法删除")
			return
		} else {
			r2 := db.Where("id = ?", val).Delete(&common.TCategory{})
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

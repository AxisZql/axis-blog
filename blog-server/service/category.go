package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Category struct{}

func (c *Category) ListCategories(*gin.Context)     {}
func (c *Category) ListCategoriesBack(*gin.Context) {}

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
func (c *Category) SaveOrUpdateCategory(*gin.Context) {}
func (c *Category) DeleteCategories(*gin.Context)     {}

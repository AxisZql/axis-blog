package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"github.com/gin-gonic/gin"
)

type Page struct{}

func (p *Page) DeletePage(ctx *gin.Context) {
	Response(ctx, errorcode.Success, nil, true, "暂时不开放删除页面功能")
}

type reqSaveOrUpdatePage struct {
	ID        int64  `json:"id" binding:"required"`
	PageCover string `json:"pageCover" binding:"required"`
	PageLabel string `json:"pageLabel" binding:"required"`
	PageName  string `json:"pageName" binding:"required"`
}

func (p *Page) SaveOrUpdatePage(ctx *gin.Context) {
	var form reqSaveOrUpdatePage
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数验证失败或无法新建页面(暂时不开放新页面创建功能)")
		return
	}
	db := common.GetGorm()
	var _page common.TPage
	r1 := db.Where("id = ?", form.ID).First(&_page)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	_page.PageName = form.PageName
	_page.PageCover = form.PageCover
	_page.PageLabel = form.PageLabel
	r1 = db.Save(&_page)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (p *Page) ListPages(ctx *gin.Context) {
	db := common.GetGorm()
	var pageList []struct {
		ID        int64  `json:"id"`
		PageCover string `json:"pageCover"`
		PageLabel string `json:"pageLabel"`
		PageName  string `json:"pageName"`
	}
	r1 := db.Model(&common.TPage{}).Find(&pageList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, pageList, true, "操作成功")
}

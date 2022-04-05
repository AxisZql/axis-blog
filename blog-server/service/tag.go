package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Tag struct{}

func (t *Tag) ListTags(*gin.Context)    {}
func (t *Tag) ListTagBack(*gin.Context) {}

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
func (t *Tag) SaveOrUpdateTag(*gin.Context) {}
func (t *Tag) DeleteTag(*gin.Context)       {}

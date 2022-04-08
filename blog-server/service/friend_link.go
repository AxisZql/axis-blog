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

type FriendLink struct{}

func (f *FriendLink) ListFriendLinks(*gin.Context) {}

type reqListFriendLinksBack struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (f *FriendLink) ListFriendLinksBack(ctx *gin.Context) {
	var form reqListFriendLinksBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	db := common.GetGorm()
	var count int64
	var linkList []struct {
		ID          int64  `json:"id"`
		LinkAddress string `json:"linkAddress"`
		LinkAvatar  string `json:"linkAvatar"`
		LinkIntro   string `json:"linkIntro"`
		LinkName    string `json:"linkName"`
	}
	if form.Keywords == "" {
		r1 := db.Model(&common.TFriendLink{}).Count(&count)
		r1 = db.Model(&common.TFriendLink{}).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&linkList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Model(&common.TFriendLink{}).Where(fmt.Sprintf("link_name LIKE %q", "%"+form.Keywords+"%")).Count(&count)
		r1 = db.Model(&common.TFriendLink{}).Where(fmt.Sprintf("link_name LIKE %q", "%"+form.Keywords+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&linkList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = linkList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqSaveOrUpdateFriendLink struct {
	ID          int64  `json:"id"`
	LinkAddress string `json:"linkAddress"`
	LinkAvatar  string `json:"linkAvatar"`
	LinkIntro   string `json:"linkIntro"`
	LinkName    string `json:"linkName"`
}

func (f *FriendLink) SaveOrUpdateFriendLink(ctx *gin.Context) {
	var form reqSaveOrUpdateFriendLink
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.ID == 0 {
		fr := common.TFriendLink{
			LinkAddress: form.LinkAddress,
			LinkAvatar:  form.LinkAvatar,
			LinkIntro:   form.LinkIntro,
			LinkName:    form.LinkName,
		}
		r1 := db.Model(&common.TFriendLink{}).Create(&fr)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		var fr common.TFriendLink
		r1 := db.Where("id = ?", form.ID).First(&fr)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		fr.LinkAvatar = form.LinkAvatar
		fr.LinkAddress = form.LinkAddress
		fr.LinkIntro = form.LinkIntro
		fr.LinkName = form.LinkName
		fr.UpdateTime = time.Now()
		r1 = db.Save(&fr)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (f *FriendLink) DeleteFriendLink(ctx *gin.Context) {
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
		r1 := db.Where("id = ?", val).Delete(&common.TFriendLink{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

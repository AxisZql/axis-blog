package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Message struct{}

type reqSaveMessage struct {
	Avatar         string `json:"avatar" binding:"required"`
	MessageContent string `json:"messageContent" binding:"required"`
	Nickname       string `json:"nickname" binding:"required"`
	Time           int    `json:"time" binding:"required"`
}

func (m *Message) SaveMessage(ctx *gin.Context) {
	var form reqSaveMessage
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	senitiveWordList := senitiveForest.GetSenitiveWord(form.MessageContent)
	if len(senitiveWordList) != 0 {
		Response(ctx, errorcode.SenitiveWordError, nil, false, fmt.Sprintf("含有敏感词:%v", senitiveWordList))
		return
	}
	db := common.GetGorm()
	message := common.TMessage{
		Avatar:         form.Avatar,
		MessageContent: form.MessageContent,
		Nickname:       form.Nickname,
		Speed:          form.Time,
	}
	r1 := db.Model(&common.TMessage{}).Create(&message)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (m *Message) ListMessage(ctx *gin.Context) {
	db := common.GetGorm()
	type ML struct {
		ID             int64  `json:"id"`
		Avatar         string `json:"avatar"`
		MessageContent string `json:"messageContent"`
		Nickname       string `json:"nickname"`
		Speed          int    `json:"time"`
	}
	var messageList []ML
	r1 := db.Model(&common.TMessage{}).Find(&messageList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, messageList, true, "操作成功")
}

type reqListMessageBack struct {
	Current  int         `form:"current"`
	Size     int         `form:"size"`
	Keywords string      `form:"keywords"`
	IsReview interface{} `form:"isReview"`
}

func (m *Message) ListMessageBack(ctx *gin.Context) {
	var form reqListMessageBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	isReview := 1
	if form.IsReview != nil {
		isReview = form.IsReview.(int)
	}
	db := common.GetGorm()
	type ML struct {
		Avatar         string    `json:"avatar"`
		CreateTime     time.Time `json:"createTime"`
		ID             int64     `json:"id"`
		IpAddress      string    `json:"ipAddress"`
		IpSource       string    `json:"ipSource"`
		IsReview       int       `json:"isReview"`
		MessageContent string    `json:"messageContent"`
		Nickname       string    `json:"nickname"`
	}
	var count int64
	var messageList []ML
	r1 := db.Model(&common.TMessage{}).Where(fmt.Sprintf("nickname LIKE %q AND is_review = ?", "%"+form.Keywords+"%"), isReview).Count(&count)
	r1 = db.Model(&common.TMessage{}).Where(fmt.Sprintf("nickname LIKE %q AND is_review = ?", "%"+form.Keywords+"%"), isReview).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&messageList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = messageList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqUpdateMessageReview struct {
	IdList   []int64 `json:"idList"`
	IsReview int     `json:"isReview"`
}

func (m *Message) UpdateMessageReview(ctx *gin.Context) {
	var form reqUpdateMessageReview
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range form.IdList {
		var t common.TMessage
		r1 := db.Where("id = ?", val).First(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if r1.Error != nil {
			continue
		}
		t.IsReview = form.IsReview
		t.UpdateTime = time.Now()
		r1 = db.Save(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (m *Message) DeleteMessage(ctx *gin.Context) {
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
		r1 := db.Where("id = ?", val).Delete(&common.TMessage{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

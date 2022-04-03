package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Talk struct{}

func (t *Talk) ListHomeTalks(ctx *gin.Context) {
	//查看最新10条说说
	db := common.GetGorm()
	talkList := make([]common.TTalk, 0)
	result := db.Model(&common.TTalk{}).Order("create_time DESC").Limit(10).Find(&talkList)
	if result.Error != nil {
		logger.Error(result.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := make([]string, 0)
	for _, val := range talkList {
		data = append(data, val.Content)
	}
	Response(ctx, errorcode.Success, data, true, "操作成功")
	return

}

type reqListTalks struct {
	Current int `form:"current" binding:"required"`
	Size    int `form:"size" binding:"required"`
}
type talk struct {
	ID           int64     `json:"id"`
	Avatar       string    `json:"avatar"`
	CommentCount int64     `json:"commentCount"`
	Content      string    `json:"content"`
	CreateTime   time.Time `json:"createTime"`
	Images       string    `json:"images"`
	ImgList      []string  `json:"imgList"`
	IsTop        int       `json:"isTop"`
	LikeCount    int64     `json:"likeCount"`
	Nickname     string    `json:"nickname"`
}

func (t *Talk) ListTalks(ctx *gin.Context) {
	var form reqListTalks
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	db := common.GetGorm()
	var talkCount int64
	var talkList []talk
	r1 := db.Model(&common.TTalk{}).Count(&talkCount)
	r2 := db.Table("v_talk_info").Limit(form.Size).Offset((form.Current - 1) * form.Size).Order("create_time DESC").Find(&talkList)
	if r1.Error != nil || r2.Error != nil {
		logger.Error(r1.Error.Error() + "|||" + r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, val := range talkList {
		var imgList []string
		if val.Images == "" {
			val.Images = "[]"
		}
		_ = json.Unmarshal([]byte(val.Images), &imgList)
		talkList[i].ImgList = imgList
	}
	data := make(map[string]interface{})
	data["count"] = talkCount
	data["recordList"] = talkList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqGetTalkById struct {
	TalkId int64 `uri:"talkId" binding:"required"`
}

func (t *Talk) GetTalkById(ctx *gin.Context) {
	var form reqGetTalkById
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var talkInfo talk
	r1 := db.Table("v_talk_info").Where("id = ?", form.TalkId).Find(&talkInfo)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var imgList []string
	if talkInfo.Images == "" {
		talkInfo.Images = "[]"
	}
	_ = json.Unmarshal([]byte(talkInfo.Images), &imgList)
	talkInfo.ImgList = imgList
	Response(ctx, errorcode.Success, talkInfo, true, "操作成功")
}
func (t *Talk) SaveTalkLike(*gin.Context)     {}
func (t *Talk) SaveTalkImages(*gin.Context)   {}
func (t *Talk) SaveOrUpdateTalk(*gin.Context) {}
func (t *Talk) DeleteTalks(*gin.Context)      {}
func (t *Talk) ListBackTalks(*gin.Context)    {}
func (t *Talk) GetBackTalkById(*gin.Context)  {}

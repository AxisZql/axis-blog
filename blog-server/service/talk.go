package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
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
	// 根据是否置顶进行排序
	sort.Slice(talkList, func(i, j int) bool {
		if talkList[i].IsTop == 1 && talkList[j].IsTop == 1 {
			return talkList[i].CreateTime.Unix() > talkList[j].CreateTime.Unix()

		} else if talkList[i].IsTop == 1 || talkList[j].IsTop == 1 {
			if talkList[i].IsTop == 1 {
				return true
			} else {
				return false
			}
		} else {
			// 两者都为0
			return talkList[i].CreateTime.Unix() > talkList[j].CreateTime.Unix()
		}
	})
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

type reqSaveTalkLike struct {
	TalkId int64  `uri:"talkId" binding:"required"`
	Path   string `uri:"like" binding:"required"`
}

func (t *Talk) SaveTalkLike(ctx *gin.Context) {
	var form reqSaveTalkLike
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Path != "like" {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	auid, ok := ctx.Get("a_userid")
	if !ok {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var ua common.TUserAuth
	r1 := db.Where("id = ?", auid).First(&ua)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var exist common.TLike
	r1 = db.Where("object = ? AND user_id = ? AND like_id = ?", "t_talk", ua.UserInfoId, form.TalkId).First(&exist)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r1.Error == nil {
		r1 = db.Model(&common.TLike{}).Where("id = ?", exist.ID).Delete(&common.TLike{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		tl := common.TLike{
			UserId: ua.UserInfoId,
			Object: "t_talk",
			LikeId: form.TalkId,
		}
		r1 = db.Model(&common.TLike{}).Create(&tl)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqSaveTalkImages struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (t *Talk) SaveTalkImages(ctx *gin.Context) {
	var form reqSaveTalkImages
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	f, _ := form.File.Open()
	extendName := strings.Split(form.File.Filename, ".")
	if len(extendName) != 2 && extendName[1] != "png" && extendName[1] != "gif" && extendName[1] != "jpg" {
		Response(ctx, errorcode.ValidError, nil, false, "不支持的图片格式;仅支持png|gif|jpg格式")
		return
	}
	defer f.Close()
	fileData, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		logger.Error(err2.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	fileMD5 := fmt.Sprintf("%x", md5.Sum(fileData))
	fileName := fileMD5 + "." + extendName[1]
	filePath := common.Conf.App.TalkDir + fileName
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	imgUrl := fmt.Sprintf("%s/talks/%s", common.Conf.App.HostName, fileName)
	Response(ctx, errorcode.Fail, imgUrl, true, "操作成功")
}

type reqSaveOrUpdateTalk struct {
	Content string `json:"content" binding:"required"`
	Images  string `json:"images"`
	IsTop   int    `json:"isTop"`
	ID      int64  `json:"id"`
	Status  int    `json:"status" binding:"required"`
}

func (t *Talk) SaveOrUpdateTalk(ctx *gin.Context) {
	var form reqSaveOrUpdateTalk
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	aid, exist := ctx.Get("a_userid")
	if !exist {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var ua common.TUserAuth
	var ui common.TUserInfo

	db := common.GetGorm()
	r := db.Where("id = ?", aid).First(&ua)
	if r.Error != nil {
		logger.Error(r.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	r = db.Where("id = ?", ua.UserInfoId).First(&ui)
	if r.Error != nil {
		logger.Error(r.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}

	if form.ID <= 0 {
		ta := common.TTalk{
			Content: form.Content,
			Images:  form.Images,
			IsTop:   form.IsTop,
			UserId:  ui.ID,
			Status:  form.Status,
		}
		r1 := db.Model(&common.TTalk{}).Create(&ta)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		var ta common.TTalk
		r1 := db.Where("id = ?", form.ID).First(&ta)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		ta.Content = form.Content
		ta.Images = form.Images
		ta.IsTop = form.IsTop
		ta.Status = form.Status
		ta.UpdateTime = time.Now()
		r1 = db.Save(&ta)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (t *Talk) DeleteTalks(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	str := fmt.Sprintf("%v", string(data))
	var idList []string
	err := json.Unmarshal([]byte(str), &idList)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range idList {
		id, _ := strconv.Atoi(val)
		r1 := db.Where("id = ?", id).Delete(&common.TTalk{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqListBackTalks struct {
	Current int `form:"current"`
	Size    int `form:"size"`
	Status  int `form:"status"`
}

func (t *Talk) ListBackTalks(ctx *gin.Context) {
	var form reqListBackTalks
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	if form.Status == 0 {
		form.Status = 1
	}
	db := common.GetGorm()
	var count int64
	r1 := db.Model(&common.TTalk{}).Where("status = ?", form.Status).Count(&count)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var talkList []struct {
		ID         int64     `json:"id"`
		Avatar     string    `json:"avatar"`
		Content    string    `json:"content"`
		CreateTime time.Time `json:"createTime"`
		Images     string    `json:"images"`
		ImgList    []string  `json:"imgList"`
		IsTop      int       `json:"isTop"`
		Nickname   string    `json:"nickname"`
		Status     int       `json:"status"`
	}
	r1 = db.Table("v_talk_info").Where("status = ?", form.Status).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&talkList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, val := range talkList {
		var tmp []string
		if val.Images == "" {
			val.Images = "[]"
		}
		err := json.Unmarshal([]byte(val.Images), &tmp)
		if err != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		talkList[i].ImgList = tmp
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = talkList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqGetBackTalkById struct {
	TalkId int64 `uri:"talkId" binding:"required"`
}

func (t *Talk) GetBackTalkById(ctx *gin.Context) {
	var form reqGetBackTalkById
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var _talk struct {
		ID         int64     `json:"id"`
		Avatar     string    `json:"avatar"`
		Content    string    `json:"content"`
		CreateTime time.Time `json:"createTime"`
		Images     string    `json:"images"`
		ImgList    []string  `json:"imgList"`
		IsTop      int       `json:"isTop"`
		Nickname   string    `json:"nickname"`
		Status     int       `json:"status"`
	}
	r1 := db.Table("v_talk_info").Where("id = ?", form.TalkId).First(&_talk)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var tmp []string
	if _talk.Images == "" {
		_talk.Images = "[]"
	}
	err := json.Unmarshal([]byte(_talk.Images), &tmp)
	if err != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	_talk.ImgList = tmp
	Response(ctx, errorcode.Success, _talk, true, "操作成功")
}

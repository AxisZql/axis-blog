package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Comment struct{}

//====分页获取评论

type reqListComment struct {
	Current int   `form:"current" binding:"required"`
	Type    int   `form:"type" biding:"required"`
	TopicId int64 `form:"topicId" biding:"required"`
}

func (c *Comment) ListComment(ctx *gin.Context) {
	var form reqListComment
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验错误")
		return
	}
	db := common.GetGorm()

	type SCL struct {
		Avatar         string    `json:"avatar"`
		CreateTime     time.Time `json:"createTime"`
		CommentContent string    `json:"commentContent"`
		ID             int64     `json:"id"`
		LikeCount      int64     `json:"likeCount"`
		Nickname       string    `json:"nickname"`
		ParentId       int64     `json:"parentId"`
		ReplyNickname  string    `json:"replyNickname"`
		ReplyUserId    int64     `json:"replyUserId"`
		ReplyWebSite   string    `json:"replyWebSite"`
		UserId         int64     `json:"userId"`
		WebSite        string    `json:"webSite"`
	}
	type CL struct {
		Avatar         string      `json:"avatar"`
		CommentContent string      `json:"commentContent"`
		CreateTime     time.Time   `json:"createTime"`
		ID             int64       `json:"id"`
		LikeCount      int64       `json:"likeCount"`
		Nickname       string      `json:"nickname"`
		ReplyCount     int64       `json:"replyCount"`
		ReplyDTOList   interface{} `json:"replyDTOList"`
		UserId         int64       `json:"userId"`
		WebSite        string      `json:"webSite"`
	}
	var commentList []CL
	var count int64
	if form.TopicId == 0 {
		r1 := db.Table("v_comment").Where("type = ? AND isNull(parent_id)", form.Type).Count(&count)
		r1 = db.Table("v_comment").Where("type = ? AND isNull(parent_id)", form.Type).Limit(10).Offset((form.Current - 1) * 10).Find(&commentList)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Table("v_comment").Where("type = ? AND topic_id = ? AND isNull(parent_id)", form.Type, form.TopicId).Count(&count)
		r1 = db.Table("v_comment").Where("type = ? AND topic_id = ? AND isNull(parent_id)", form.Type, form.TopicId).Limit(10).Offset((form.Current - 1) * 10).Find(&commentList)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	for i, val := range commentList {
		var _count int64
		r2 := db.Model(&common.TComment{}).Where("parent_id = ?", val.ID).Count(&_count)
		var scl []SCL
		r2 = db.Table("v_comment").Where("parent_id = ?", val.ID).Limit(5).Find(&scl)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		commentList[i].ReplyCount = _count
		likeCount, err := common.GetCommentLikeCountById(val.ID)
		if err != nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for k, v := range scl {
			lc, e := common.GetCommentLikeCountById(v.ID)
			if e != nil {
				logger.Error(e.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			scl[k].LikeCount = lc
		}
		commentList[i].LikeCount = likeCount
		commentList[i].ReplyDTOList = scl
	}
	//data := common.ConvertCommentData(commentList)
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = commentList
	Response(ctx, errorcode.Success, data, true, "操作成功")

}

type reqSaveComment struct {
	CommentContent string      `json:"commentContent" binding:"required"`
	ParentId       interface{} `json:"parentId"`
	ReplyUserId    interface{} `json:"replyUserId"`
	TopicId        interface{} `json:"topicId"`
	Type           int         `json:"type" binding:"required"`
}

func (c *Comment) SaveComment(ctx *gin.Context) {
	var form reqSaveComment
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	senitiveWordList := senitiveForest.GetSenitiveWord(form.CommentContent)
	if len(senitiveWordList) != 0 {
		Response(ctx, errorcode.SenitiveWordError, nil, false, fmt.Sprintf("含有敏感词:%v", senitiveWordList))
		return
	}
	db := common.GetGorm()
	userid, exist := ctx.Get("a_userid")
	if !exist {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var ua common.TUserAuth
	r := db.Where("id = ?", userid).First(&ua)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		logger.Error(r.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	filed := make([]string, 0)
	comment := common.TComment{
		CommentContent: form.CommentContent,
		Type:           form.Type,
		UserId:         ua.UserInfoId,
	}
	filed = append(filed, []string{"comment_content", "type", "user_id"}...)
	if form.TopicId != nil {
		filed = append(filed, "topic_id")
		tid, _ := strconv.Atoi(form.TopicId.(string))
		comment.TopicId = int64(tid)
	}
	if form.ParentId != nil && form.ReplyUserId != nil {
		filed = append(filed, []string{"parent_id", "reply_user_id"}...)
		comment.ParentId = int64(form.ParentId.(float64))
		comment.ReplyUserId = int64(form.ReplyUserId.(float64))
	}

	r1 := db.Select(filed).Create(&comment)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqListRepliesByCommentId struct {
	Current   int    `form:"current"`
	Size      int    `form:"size"`
	CommentId int64  `uri:"commentId" binding:"required"`
	Path      string `uri:"replies" binding:"required"`
}

func (c *Comment) ListRepliesByCommentId(ctx *gin.Context) {
	var form reqListRepliesByCommentId
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	_ = ctx.ShouldBind(&form)
	if form.Path != "replies" {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 5
	}
	db := common.GetGorm()
	type CL struct {
		Avatar         string    `json:"avatar"`
		CommentContent string    `json:"commentContent"`
		CreateTime     time.Time `json:"createTime"`
		ID             int64     `json:"id"`
		LikeCount      int64     `json:"likeCount"`
		Nickname       string    `json:"nickname"`
		ParentId       int64     `json:"parentId"`
		ReplyNickname  string    `json:"replyNickname"`
		ReplyUserId    int64     `json:"replyUserId"`
		ReplyWebSite   string    `json:"replyWebSite"`
		UserId         int64     `json:"userId"`
		WebSite        string    `json:"webSite"`
	}
	var commentList []CL
	r1 := db.Table("v_comment").Where("parent_id = ?", form.CommentId).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&commentList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, val := range commentList {
		var count int64
		r2 := db.Model(&common.TLike{}).Where("object = ? and like_id = ?", "t_comment", val.ID).Count(&count)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		commentList[i].LikeCount = count
	}
	Response(ctx, errorcode.Success, commentList, true, "操作成功")
}

type reqSaveCommentLike struct {
	CommentId int64  `uri:"commentId" binding:"required"`
	Path      string `uri:"like" binding:"required"`
}

func (c *Comment) SaveCommentLike(ctx *gin.Context) {
	var form reqSaveCommentLike
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
	r1 = db.Where("object = ? AND user_id = ? AND like_id = ?", "t_comment", ua.UserInfoId, form.CommentId).First(&exist)
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
			Object: "t_comment",
			LikeId: form.CommentId,
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

type reqUpdateCommentReview struct {
	IdList   []int64 `json:"idList"`
	IsReview int     `json:"isReview"`
}

func (c *Comment) UpdateCommentReview(ctx *gin.Context) {
	var form reqUpdateCommentReview
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range form.IdList {
		var t common.TComment
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
func (c *Comment) DeleteComment(ctx *gin.Context) {
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
		var chile []common.TComment
		r1 := db.Where("parent_id = ?", val).Find(&chile)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		var t common.TComment
		r1 = db.Where("id = ?", val).First(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for _, c := range chile {
			c.IsDelete = 1
			c.UpdateTime = time.Now()
			r2 := db.Save(&c)
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		t.IsDelete = 1
		t.UpdateTime = time.Now()
		r1 = db.Save(&t)
		r1 = db.Where("id = ?", val).First(&t)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqListCommentBack struct {
	Current  int         `form:"current"`
	Size     int         `form:"size"`
	IsReview interface{} `form:"isReview"`
	Keywords string      `form:"keywords"`
	Type     int         `form:"type"`
}

func (c *Comment) ListCommentBack(ctx *gin.Context) {
	var form reqListCommentBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Type == 0 {
		form.Type = 1
	}
	isReview := 1
	if form.IsReview != nil {
		isReview = form.IsReview.(int)
	}
	db := common.GetGorm()
	type CI struct {
		ArticleTitle   string    `json:"articleTitle"`
		Avatar         string    `json:"avatar"`
		CommentContent string    `json:"commentContent"`
		CreatTime      time.Time `json:"creatTime"`
		ID             int64     `json:"id"`
		IsReview       int       `json:"isReview"`
		Nickname       string    `json:"nickname"`
		ReplyNickname  string    `json:"replyNickname"`
		Type           int       `json:"type"`
		TopicId        int64     `json:"topicId"`
	}
	var count int64
	var commentList []CI

	r1 := db.Table("v_comment").Where(fmt.Sprintf("type = ? AND is_review = ? AND (nickname LIKE %q OR reply_nickname LIKE %q)", "%"+form.Keywords+"%", "%"+form.Keywords+"%"), form.Type, isReview).Count(&count)
	r1 = db.Table("v_comment").Where(fmt.Sprintf("type = ? AND is_review = ? AND (nickname LIKE %q OR reply_nickname LIKE %q)", "%"+form.Keywords+"%", "%"+form.Keywords+"%"), form.Type, isReview).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&commentList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if form.Type == 1 {
		for i, val := range commentList {
			var t common.TArticle
			r2 := db.Where("id = ?", val.TopicId).First(&t)
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			commentList[i].ArticleTitle = t.ArticleTitle
		}
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = commentList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

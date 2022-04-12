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
	offset := (form.Current - 1) * 10
	rows, err := db.Raw("select * from v_comment where type = ? and topic_id = ? limit ?,?", form.Type, form.TopicId, offset, 10).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	commentList := make([]common.VComment, 0)
	for rows.Next() {
		var t common.VComment
		_ = db.ScanRows(rows, &t)
		commentList = append(commentList, t)
	}
	data := common.ConvertCommentData(commentList)
	Response(ctx, errorcode.Success, data, true, "操作成功")

}
func (c *Comment) SaveComment(*gin.Context)            {}
func (c *Comment) ListRepliesByCommentId(*gin.Context) {}
func (c *Comment) SaveCommentLike(*gin.Context)        {}

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

package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Comment struct {
	ctrl.CommentHandle
}

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
func (c *Comment) UpdateCommentReview(*gin.Context)    {}
func (c *Comment) DeleteComment(*gin.Context)          {}
func (c *Comment) ListCommentBack(*gin.Context)        {}

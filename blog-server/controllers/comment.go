package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 9:08 PM
* @desc: 评论模块接口
 */

type CommentHandle interface {
	ListComment(*gin.Context)            //查询评论
	SaveComment(*gin.Context)            //添加评论
	ListRepliesByCommentId(*gin.Context) //查询评论下的回复
	SaveCommentLike(*gin.Context)        //评论点赞
	UpdateCommentReview(*gin.Context)    //审核评论
	DeleteComment(*gin.Context)          //删除评论
	ListCommentBack(*gin.Context)        //查询后台评论列表
}

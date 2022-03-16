package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Comment struct {
	ctrl.CommentHandle
}

func (c *Comment) ListComment(*gin.Context)            {}
func (c *Comment) SaveComment(*gin.Context)            {}
func (c *Comment) ListRepliesByCommentId(*gin.Context) {}
func (c *Comment) SaveCommentLike(*gin.Context)        {}
func (c *Comment) UpdateCommentReview(*gin.Context)    {}
func (c *Comment) DeleteComment(*gin.Context)          {}
func (c *Comment) ListCommentBack(*gin.Context)        {}

package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Tag struct {
	ctrl.TagHandle
}

func (t *Tag) ListTags(*gin.Context)        {}
func (t *Tag) ListTagBack(*gin.Context)     {}
func (t *Tag) ListTagBySearch(*gin.Context) {}
func (t *Tag) SaveOrUpdateTag(*gin.Context) {}
func (t *Tag) DeleteTag(*gin.Context)       {}

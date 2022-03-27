package service

import (
	"github.com/gin-gonic/gin"
)

type Tag struct{}

func (t *Tag) ListTags(*gin.Context)        {}
func (t *Tag) ListTagBack(*gin.Context)     {}
func (t *Tag) ListTagBySearch(*gin.Context) {}
func (t *Tag) SaveOrUpdateTag(*gin.Context) {}
func (t *Tag) DeleteTag(*gin.Context)       {}

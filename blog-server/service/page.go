package service

import (
	"github.com/gin-gonic/gin"
)

type Page struct{}

func (p *Page) DeletePage(*gin.Context)       {}
func (p *Page) SaveOrUpdatePage(*gin.Context) {}
func (p *Page) ListPages(*gin.Context)        {}

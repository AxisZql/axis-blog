package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Page struct {
	ctrl.PageHandle
}

func (p *Page) DeletePage(*gin.Context)       {}
func (p *Page) SaveOrUpdatePage(*gin.Context) {}
func (p *Page) ListPages(*gin.Context)        {}

package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Menu struct {
	ctrl.MenuHandle
}

func (m *Menu) ListMenus(*gin.Context)        {}
func (m *Menu) SaveOrUpdateMenu(*gin.Context) {}
func (m *Menu) DeleteMenu(*gin.Context)       {}
func (m *Menu) ListMenuOptions(*gin.Context)  {}
func (m *Menu) ListUserMenus(*gin.Context)    {}

package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Category struct {
	ctrl.CategoryHandle
}

func (c *Category) ListCategories(*gin.Context)         {}
func (c *Category) ListCategoriesBack(*gin.Context)     {}
func (c *Category) ListCategoriesBySearch(*gin.Context) {}
func (c *Category) SaveOrUpdateCategory(*gin.Context)   {}
func (c *Category) DeleteCategories(*gin.Context)       {}

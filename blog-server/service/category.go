package service

import (
	"github.com/gin-gonic/gin"
)

type Category struct{}

func (c *Category) ListCategories(*gin.Context)         {}
func (c *Category) ListCategoriesBack(*gin.Context)     {}
func (c *Category) ListCategoriesBySearch(*gin.Context) {}
func (c *Category) SaveOrUpdateCategory(*gin.Context)   {}
func (c *Category) DeleteCategories(*gin.Context)       {}

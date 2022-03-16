package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Resource struct {
	ctrl.ResourceHandle
}

func (r *Resource) ListResources(*gin.Context)        {}
func (r *Resource) DeleteResource(*gin.Context)       {}
func (r *Resource) SaveOrUpdateResource(*gin.Context) {}
func (r *Resource) ListResourceOption(*gin.Context)   {}

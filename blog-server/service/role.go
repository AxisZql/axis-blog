package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Role struct {
	ctrl.RoleHandle
}

func (r *Role) ListUserRoles(*gin.Context)    {}
func (r *Role) ListRoles(*gin.Context)        {}
func (r *Role) SaveOrUpdateRole(*gin.Context) {}
func (r *Role) DeleteRoles(*gin.Context)      {}

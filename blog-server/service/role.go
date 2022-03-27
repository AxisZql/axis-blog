package service

import (
	"github.com/gin-gonic/gin"
)

type Role struct{}

func (r *Role) ListUserRoles(*gin.Context)    {}
func (r *Role) ListRoles(*gin.Context)        {}
func (r *Role) SaveOrUpdateRole(*gin.Context) {}
func (r *Role) DeleteRoles(*gin.Context)      {}

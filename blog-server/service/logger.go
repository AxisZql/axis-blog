package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Logger struct {
	ctrl.LoggerHandler
}

func (l *Logger) ListOperationLogs(*gin.Context)   {}
func (l *Logger) DeleteOperationLogs(*gin.Context) {}

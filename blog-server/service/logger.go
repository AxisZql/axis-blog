package service

import (
	"github.com/gin-gonic/gin"
)

type Logger struct{}

func (l *Logger) ListOperationLogs(*gin.Context)   {}
func (l *Logger) DeleteOperationLogs(*gin.Context) {}

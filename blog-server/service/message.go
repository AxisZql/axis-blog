package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Message struct {
	ctrl.MessageHandle
}

func (m *Message) SaveMessage(*gin.Context)         {}
func (m *Message) ListMessage(*gin.Context)         {}
func (m *Message) ListMessageBack(*gin.Context)     {}
func (m *Message) UpdateMessageReview(*gin.Context) {}
func (m *Message) DeleteMessage(*gin.Context)       {}

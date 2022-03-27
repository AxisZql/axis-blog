package service

import (
	"github.com/gin-gonic/gin"
)

type Message struct{}

func (m *Message) SaveMessage(*gin.Context)         {}
func (m *Message) ListMessage(*gin.Context)         {}
func (m *Message) ListMessageBack(*gin.Context)     {}
func (m *Message) UpdateMessageReview(*gin.Context) {}
func (m *Message) DeleteMessage(*gin.Context)       {}

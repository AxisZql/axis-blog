package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Talk struct {
	ctrl.TalkHandle
}

func (t *Talk) ListHomeTalks(*gin.Context)    {}
func (t *Talk) ListTalks(*gin.Context)        {}
func (t *Talk) GetTalkById(*gin.Context)      {}
func (t *Talk) SaveTalkLike(*gin.Context)     {}
func (t *Talk) SaveTalkImages(*gin.Context)   {}
func (t *Talk) SaveOrUpdateTalk(*gin.Context) {}
func (t *Talk) DeleteTalks(*gin.Context)      {}
func (t *Talk) ListBackTalks(*gin.Context)    {}
func (t *Talk) GetBackTalkById(*gin.Context)  {}

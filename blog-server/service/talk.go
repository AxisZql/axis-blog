package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Talk struct {
	ctrl.TalkHandle
}

func (t *Talk) ListHomeTalks(ctx *gin.Context) {
	//查看最新10条说说
	db := common.GetGorm()
	talkList := make([]common.TTalk, 0)
	result := db.Model(&common.TTalk{}).Order("create_time DESC").Limit(10).Find(&talkList)
	if result.Error != nil {
		logger.Error(result.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := make([]string, 0)
	for _, val := range talkList {
		data = append(data, val.Content)
	}
	Response(ctx, errorcode.Success, data, true, "操作成功")
	return

}
func (t *Talk) ListTalks(*gin.Context)        {}
func (t *Talk) GetTalkById(*gin.Context)      {}
func (t *Talk) SaveTalkLike(*gin.Context)     {}
func (t *Talk) SaveTalkImages(*gin.Context)   {}
func (t *Talk) SaveOrUpdateTalk(*gin.Context) {}
func (t *Talk) DeleteTalks(*gin.Context)      {}
func (t *Talk) ListBackTalks(*gin.Context)    {}
func (t *Talk) GetBackTalkById(*gin.Context)  {}

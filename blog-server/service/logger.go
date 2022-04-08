package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

type Logger struct{}

type reqListOperationLogs struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (l *Logger) ListOperationLogs(ctx *gin.Context) {
	var form reqListOperationLogs
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	db := common.GetGorm()
	var count int64
	var logList []struct {
		CreateTime    time.Time `json:"createTime"`
		ID            int64     `json:"id"`
		IpAddress     string    `json:"ipAddress"`
		IpSource      string    `json:"ipSource"`
		Nickname      string    `json:"nickname"`
		OptDesc       string    `json:"optDesc"`
		OptMethod     string    `json:"optMethod"`
		OptModule     string    `json:"optModule"`
		OptType       string    `json:"optType"`
		OptUrl        string    `json:"optUrl"`
		RequestMethod string    `json:"requestMethod"`
		RequestParam  string    `json:"requestParam"`
		ResponseData  string    `json:"responseData"`
	}
	if form.Keywords == "" {
		r1 := db.Model(&common.TOperationLog{}).Count(&count)
		r1 = db.Model(&common.TOperationLog{}).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&logList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		k := form.Keywords
		r1 := db.Model(&common.TOperationLog{}).Where(fmt.Sprintf("ip_address = %q OR opt_module LIKE %q OR ip_source LIKE %q OR nickname LIKE %q", k, "%"+k+"%", "%"+k+"%", "%"+k+"%")).Count(&count)
		r1 = db.Model(&common.TOperationLog{}).Where(fmt.Sprintf("ip_address = %q OR opt_module LIKE %q OR ip_source LIKE %q OR nickname LIKE %q", k, "%"+k+"%", "%"+k+"%", "%"+k+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&logList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = logList
	Response(ctx, errorcode.Success, data, true, "操作成功")

}
func (l *Logger) DeleteOperationLogs(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	str := fmt.Sprintf("%v", string(data))
	var idList []int64
	err := json.Unmarshal([]byte(str), &idList)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	db := common.GetGorm()
	for _, val := range idList {
		r1 := db.Where("id = ?", val).Delete(&common.TOperationLog{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

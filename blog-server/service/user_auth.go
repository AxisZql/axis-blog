package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
)

type UserAuth struct {
	ctrl.UserAuthHandle
}

func (user *UserAuth) SendEmailCode(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"a": "fuck you"})
}

type reqListUser struct {
	Type int `form:"type" binding:"required"`
}

type uArea struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (user *UserAuth) ListUserAreas(ctx *gin.Context) {
	var form reqListUser
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	switch form.Type {
	//获取用户地区分布
	case 1:
		db := common.GetGorm()
		uAreaList := make([]uArea, 0)
		rows, err := db.Raw("select ip_source as name,count(*) as value from t_user_auth group by ip_source;").Rows()
		if err != nil {
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for rows.Next() {
			var t uArea
			_ = db.ScanRows(rows, &t)
			uAreaList = append(uAreaList, t)
		}
		Response(ctx, errorcode.Success, uAreaList, true, "操作成功")
		return
	// 获取游客地区分布
	case 2:
		redisClient := common.GetRedis()
		u, err := redisClient.HGetAll(rediskey.VisitorArea).Result()
		if err != nil && err != redis.Nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		uAreaList := make([]uArea, 0)
		for key, val := range u {
			v, _ := strconv.Atoi(val)
			uAreaList = append(uAreaList, uArea{Name: key, Value: int64(v)})
		}
		Response(ctx, errorcode.Success, uAreaList, true, "操作成功")
		return
	}
	Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
}

func (user *UserAuth) ListUsers(ctx *gin.Context) {

}
func (user *UserAuth) Register(*gin.Context)            {}
func (user *UserAuth) UpdatePassword(*gin.Context)      {}
func (user *UserAuth) UpdateAdminPassword(*gin.Context) {}
func (user *UserAuth) WeiboLogin(*gin.Context)          {}
func (user *UserAuth) QQLogin(*gin.Context)             {}

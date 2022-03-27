package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

type UserAuth struct{}

//====== 通过QQ邮箱给新用户发送验证码

type reqSendEmailCode struct {
	Username string `form:"username" binding:"required"`
}

func (user *UserAuth) SendEmailCode(ctx *gin.Context) {
	var form reqSendEmailCode
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验错误")
		return
	}
	redisClient := common.GetRedis()
	// 生成6位验证码
	code := common.GetRandStr(6)
	// 在redis中写入验证码
	err := redisClient.Set(fmt.Sprintf(rediskey.UserCodeKey, form.Username), code, rediskey.CodeExpireTime).Err()
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	err = common.SentCodeByEmail(code, form.Username)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")

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

//========新用户注册
type reqRegister struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (user *UserAuth) Register(ctx *gin.Context) {
	var form reqRegister
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	redisClient := common.GetRedis()
	code, err := redisClient.Get(fmt.Sprintf(rediskey.UserCodeKey, form.Username)).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if err == redis.Nil || strings.ToLower(code) != strings.ToLower(form.Code) {
		Response(ctx, errorcode.ValidError, nil, false, "验证码错误")
		return
	}
	//查看对应账号是否存在
	db := common.GetGorm()
	pr1 := db.Where("username = ?", form.Username).First(&common.TUserAuth{})
	pr2 := db.Where("email = ?", form.Username).First(&common.TUserInfo{})
	if pr1.Error == nil || pr2.Error == nil {
		Response(ctx, errorcode.UsernameExistError, nil, false, "用户名已存在")
		return
	}

	// 获取网站配置
	config, err := redisClient.Get(rediskey.WebsiteConfig).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	} else if err == redis.Nil {
		w := common.TWebsiteConfig{}
		r := db.Model(&common.TWebsiteConfig{}).First(&w)
		if r.Error != nil {
			logger.Error(r.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		config = w.Config
		if err := redisClient.Set(rediskey.WebsiteConfig, config, -1).Err(); err != nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	wConfig := webConfig{}
	if err := json.Unmarshal([]byte(config), &wConfig); err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 新建用户信息
	userInfo := common.TUserInfo{
		Email:    form.Username,
		Nickname: fmt.Sprintf("用户%d", time.Now().Unix()),
		Avatar:   wConfig.UserAvatar,
	}
	r2 := db.Create(&userInfo)
	if r2.Error != nil {
		logger.Error(r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	pwd, err := common.EncryptionPwd(form.Password)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 把对应账号信息添加到t_user_auth表中
	ip := ctx.ClientIP()
	ipInfo, _ := common.GetIpAddressAndSource(ip)
	var ipSource string
	if ipInfo != nil {
		ipSource = ipInfo.Data.Province
	}
	userAuth := common.TUserAuth{
		UserInfoId:    userInfo.ID,
		Username:      form.Username,
		Password:      pwd,
		LoginType:     1,
		IpAddress:     ip,
		IpSource:      ipSource,
		UserAgent:     ctx.GetHeader("User-Agent"),
		LastLoginTime: time.Now(),
	}
	r3 := db.Create(&userAuth)
	if r3.Error != nil {
		logger.Error(r3.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")

}
func (user *UserAuth) UpdatePassword(*gin.Context)      {}
func (user *UserAuth) UpdateAdminPassword(*gin.Context) {}
func (user *UserAuth) WeiboLogin(*gin.Context)          {}
func (user *UserAuth) QQLogin(*gin.Context)             {}

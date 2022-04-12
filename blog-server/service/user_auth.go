package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
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

type reqListUsers struct {
	Current   int    `form:"current"`
	Size      int    `form:"size"`
	Keywords  string `form:"keywords"`
	LoginType int    `form:"loginType"`
}

func (user *UserAuth) ListUsers(ctx *gin.Context) {
	var form reqListUsers
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	if form.LoginType == 0 {
		form.LoginType = 1
	}
	var count int64
	type RL struct {
		ID       int64  `json:"id"`
		RoleName string `json:"roleName"`
	}
	var userList []struct {
		Avatar        string    `json:"avatar"`
		CreateTime    time.Time `json:"createTime"`
		ID            int64     `json:"id"`
		IpAddress     string    `json:"ipAddress"`
		IpSource      string    `json:"ipSource"`
		IsDisable     int       `json:"isDisable"`
		LastLoginTime time.Time `json:"lastLoginTime"`
		Nickname      string    `json:"nickname"`
		Status        int       `json:"status"`
		UserInfoId    int64     `json:"userInfoId"`
		RoleList      []RL      `json:"roleList"`
	}
	if form.Keywords == "" {
		r1 := db.Table("v_user_info").Where("login_type = ?", form.LoginType).Count(&count)
		r1 = db.Table("v_user_info").Where("login_type = ?", form.LoginType).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&userList)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		r1 := db.Table("v_user_info").Where(fmt.Sprintf("login_type = %d AND nickname LIKE %q", form.LoginType, "%"+form.Keywords+"%")).Count(&count)
		r1 = db.Table("v_user_info").Where(fmt.Sprintf("login_type = %d AND nickname LIKE %q", form.LoginType, "%"+form.Keywords+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&userList)
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	for i, u := range userList {
		var tur []common.TUserRole
		roleList := make([]RL, 0)
		r2 := db.Model(&common.TUserRole{}).Where("user_id = ?", u.UserInfoId).Find(&tur)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for _, val := range tur {
			var role RL
			r3 := db.Model(&common.TRole{}).Where("id = ?", val.RoleId).Find(&role)
			if r3.Error != nil && r3.Error != gorm.ErrRecordNotFound {
				logger.Error(r3.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			roleList = append(roleList, role)
		}
		userList[i].RoleList = roleList
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = userList
	Response(ctx, errorcode.Success, data, true, "操作成功")
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
	var role common.TRole
	r4 := db.Where("role_label = ?", "user").First(&role)
	if r4.Error != nil && r4.Error != gorm.ErrRecordNotFound {
		logger.Error(r4.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r4.Error == gorm.ErrRecordNotFound {
		role.RoleName = "用户"
		role.RoleLabel = "user"
		role.IsDisable = 0
		r4 = db.Model(&common.TRole{}).Create(&role)
		if r4.Error != nil {
			logger.Error(r4.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	userRole := common.TUserRole{
		UserId: userInfo.ID,
		RoleId: role.ID,
	}
	r5 := db.Model(&common.TUserRole{}).Create(&userRole)
	if r5.Error != nil {
		logger.Error(r5.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")

}
func (user *UserAuth) UpdatePassword(*gin.Context) {}

type reqUpdateAdminPassword struct {
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	OldPassword     string `json:"oldPassword" binding:"required"`
}

func (user *UserAuth) UpdateAdminPassword(ctx *gin.Context) {
	var form reqUpdateAdminPassword
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	_session, _ := Store.Get(ctx.Request, "CurUser")
	auid := _session.Values["a_userid"]
	var au common.TUserAuth
	r1 := db.Where("id = ?", auid).First(&au)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if !common.VerifyPwd(au.Password, form.OldPassword) {
		Response(ctx, errorcode.AuthorizedError, nil, false, "密码错误")
		return
	}
	if form.NewPassword != form.ConfirmPassword {
		Response(ctx, errorcode.ValidError, nil, false, "两次密码输入不一致")
		return
	}
	pwd, err := common.EncryptionPwd(form.NewPassword)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	au.Password = pwd
	au.UpdateTime = time.Now()
	r1 = db.Save(&au)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (user *UserAuth) WeiboLogin(*gin.Context) {}
func (user *UserAuth) QQLogin(*gin.Context)    {}

package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	ctrl "blog-server/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserInfo struct {
	ctrl.UserInfoHandle
}

func (u *UserInfo) UpdateUserInfo(*gin.Context)   {}
func (u *UserInfo) UpdateUserAvatar(*gin.Context) {}

//=====用户换绑邮箱
type reqSaveUserEmail struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (u *UserInfo) SaveUserEmail(ctx *gin.Context) {
	var form reqSaveUserEmail
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	redisClient := common.GetRedis()
	code, err := redisClient.Get(fmt.Sprintf(rediskey.UserCodeKey, form.Email)).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if err == redis.Nil || strings.ToLower(code) != strings.ToLower(form.Code) {
		Response(ctx, errorcode.ValidError, nil, false, "验证码错误")
		return
	}
	//查看对应邮箱是否被绑定
	db := common.GetGorm()
	pr1 := db.Where("username = ?", form.Email).First(&common.TUserAuth{})
	pr2 := db.Where("email = ?", form.Email).First(&common.TUserInfo{})
	if pr1.Error == nil || pr2.Error == nil {
		Response(ctx, errorcode.UsernameExistError, nil, false, "该邮箱已经被绑定")
		return
	}
	_session, _ := Store.Get(ctx.Request, "CurUser")
	// 获取当前用户id
	userid := _session.Values["a_userid"]
	userInfo := common.TUserInfo{}
	userAuth := common.TUserAuth{}
	r1 := db.Where("id = ?", userid).First(&userAuth)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r1.Error != nil {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "用户名不存在")
		return
	}
	r2 := db.Where("id = ?", userAuth.UserInfoId).First(&userInfo)
	if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r2.Error != nil {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "用户名不存在")
		return
	}
	userInfo.Email = form.Email
	userInfo.UpdateTime = time.Now()
	userAuth.Username = form.Email
	userAuth.UpdateTime = time.Now()
	r3 := db.Save(&userInfo)
	r4 := db.Save(&userAuth)
	if r3.Error != nil || r4.Error != nil {
		logger.Error(r3.Error.Error() + "|||" + r4.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return

	}
	// 清除登陆状态
	for key, _ := range _session.Values {
		delete(_session.Values, key)
	}
	_ = _session.Save(ctx.Request, ctx.Writer)
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (u *UserInfo) UpdateUserRole(*gin.Context)    {}
func (u *UserInfo) UpdateUserDisable(*gin.Context) {}
func (u *UserInfo) ListOnlineUsers(*gin.Context)   {}
func (u *UserInfo) RemoveOnlineUser(*gin.Context)  {}

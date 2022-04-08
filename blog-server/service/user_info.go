package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type UserInfo struct{}

type reqUpdateUserInfo struct {
	Intro     string `json:"intro" binding:"required"`
	LoginType int    `json:"loginType"`
	Nickname  string `json:"nickname" binding:"required"`
	WebSite   string `json:"webSite" binding:"required"`
}

func (u *UserInfo) UpdateUserInfo(ctx *gin.Context) {
	var form reqUpdateUserInfo
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
	}
	db := common.GetGorm()
	_session, _ := Store.Get(ctx.Request, "CurUser")
	userid := _session.Values["a_userid"]
	userAuth := common.TUserAuth{}
	userInfo := common.TUserInfo{}
	r1 := db.Where("id = ?", userid).First(&userAuth)
	r2 := db.Where("id = ?", userAuth.UserInfoId).First(&userInfo)
	if r1.Error != nil || r2.Error != nil {
		if r1.Error == gorm.ErrRecordNotFound || r2.Error == gorm.ErrRecordNotFound {
			Response(ctx, errorcode.UsernameNotExist, nil, false, "用户不存在")
			return
		}
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	userInfo.Intro = form.Intro
	userInfo.Nickname = form.Nickname
	userInfo.WebSite = form.WebSite
	userInfo.UpdateTime = time.Now()
	r3 := db.Save(&userInfo)
	if r3.Error != nil {
		logger.Error(r3.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqUpdateUserAvatar struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (u *UserInfo) UpdateUserAvatar(ctx *gin.Context) {
	var form reqUpdateUserAvatar
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, true, "参数校验失败")
		return
	}
	_session, _ := Store.Get(ctx.Request, "CurUser")
	userid := _session.Values["a_userid"]
	db := common.GetGorm()
	ua := common.TUserAuth{}
	ui := common.TUserInfo{}
	r1 := db.Where("id = ?", userid).First(&ua)
	r2 := db.Where("id = ?", ua.UserInfoId).First(&ui)
	if r1.Error != nil || r2.Error != nil {
		if r1.Error == gorm.ErrRecordNotFound || r2.Error == gorm.ErrRecordNotFound {
			Response(ctx, errorcode.UsernameNotExist, nil, false, "用户不存在")
			return
		}
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	f, _ := form.File.Open()
	extendName := strings.Split(form.File.Filename, ".")
	if len(extendName) != 2 && extendName[1] != "png" && extendName[1] != "gif" && extendName[1] != "jpg" {
		Response(ctx, errorcode.ValidError, nil, false, "不支持的图片格式;仅支持png|gif|jpg格式")
		return
	}
	defer f.Close()
	fileData, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		logger.Error(err2.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	preFilePath := common.Conf.App.AvatarDir + strings.Split(ui.Avatar, "avatar/")[1]
	fileNameMD5 := fmt.Sprintf("%x", md5.Sum(fileData))
	filePath := common.Conf.App.AvatarDir + fileNameMD5 + "." + extendName[1]
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	ui.Avatar = fmt.Sprintf("%s:%d/avatar/%s.%s", common.Conf.App.HostName, common.Conf.App.Port, fileNameMD5, extendName[1])
	ui.UpdateTime = time.Now()
	r3 := db.Save(&ui)
	if r3.Error != nil {
		logger.Error(r3.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if filePath != preFilePath {
		err = os.Remove(preFilePath)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	Response(ctx, errorcode.Success, ui.Avatar, true, "操作成功")
}

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

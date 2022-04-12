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
	"strconv"
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

type reqUpdateUserRole struct {
	Nickname   string  `json:"nickname" binding:"required"`
	UserInfoId int64   `json:"userInfoId" binding:"required"`
	RoleIdList []int64 `json:"roleIdList" binding:"required"`
	ID         int64   `json:"id" binding:"required"`
}

func (u *UserInfo) UpdateUserRole(ctx *gin.Context) {
	var form reqUpdateUserRole
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var user common.TUserInfo
	r1 := db.Model(&common.TUserInfo{}).Where("id = ?", form.UserInfoId).First(&user)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r1.Error == gorm.ErrRecordNotFound {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "用户不存在")
		return
	}
	user.Nickname = form.Nickname
	user.UpdateTime = time.Now()
	r1 = db.Save(&user)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for _, val := range form.RoleIdList {
		r2 := db.Where("role_id = ? AND user_id = ?", val, form.UserInfoId).First(&common.TUserRole{})
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if r2.Error == nil {
			continue
		}
		t := common.TUserRole{
			RoleId: val,
			UserId: form.UserInfoId,
		}
		r3 := db.Model(&common.TUserRole{}).Create(&t)
		if r3.Error != nil {
			logger.Error(r3.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	// 删除不属于该用户的权限标签
	canUse := make(map[int64]bool)
	for _, val := range form.RoleIdList {
		canUse[val] = true
	}
	var turList []common.TUserRole
	r4 := db.Where("user_id = ?", form.UserInfoId).Find(&turList)
	if r4.Error != nil && r4.Error != gorm.ErrRecordNotFound {
		logger.Error(r4.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for _, val := range turList {
		if _, ok := canUse[val.RoleId]; ok {
			continue
		}
		r5 := db.Where("id = ?", val.ID).Delete(&common.TUserRole{})
		if r5.Error != nil && r5.Error != gorm.ErrRecordNotFound {
			logger.Error(r5.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqUpdateUserDisable struct {
	ID        int64 `json:"id"`
	IsDisable int   `json:"isDisable"`
}

func (u *UserInfo) UpdateUserDisable(ctx *gin.Context) {
	var form reqUpdateUserDisable
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var ui common.TUserInfo
	r1 := db.Where("id = ?", form.ID).First(&ui)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	ui.IsDisable = form.IsDisable
	ui.UpdateTime = time.Now()
	r1 = db.Save(&ui)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqListOnlineUsers struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (u *UserInfo) ListOnlineUsers(ctx *gin.Context) {
	var form reqListOnlineUsers
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	type UF struct {
		Avatar        string    `json:"avatar"`
		Browser       string    `json:"browser"`
		IpAddress     string    `json:"ipAddress"`
		IpSource      string    `json:"ipSource"`
		LastLoginTime time.Time `json:"lastLoginTime"`
		Nickname      string    `json:"nickname"`
		Os            string    `json:"os"`
		UserInfoId    int64     `json:"userInfoId"`
	}
	redisClient := common.GetRedis()
	us, err := redisClient.SMembers(rediskey.OnlineUser).Result()
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	//count := len(us)
	var ufList []UF
	uIdList := make([]int64, 0)
	for _, val := range us {
		t, _ := strconv.Atoi(val)
		uIdList = append(uIdList, int64(t))
	}
	db := common.GetGorm()
	r1 := db.Table("v_user_info").Where(fmt.Sprintf("id IN ? AND nickname LIKE %q", "%"+form.Keywords+"%"), uIdList).Find(&ufList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := make(map[string]interface{})
	data["count"] = len(ufList)
	data["recordList"] = ufList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqRemoveOnlineUser struct {
	UserInfoId int64  `uri:"userInfoId" binding:"required"`
	Path       string `uri:"online" binding:"required"`
}

func (u *UserInfo) RemoveOnlineUser(ctx *gin.Context) {
	var form reqRemoveOnlineUser
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Path != "online" {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	redisClient := common.GetRedis()
	err := redisClient.SRem(rediskey.OnlineUser, form.UserInfoId).Err()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

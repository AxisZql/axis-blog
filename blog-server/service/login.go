package service

import (
	"blog-server/common"
	"blog-server/common/auth"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mssola/user_agent"
	"gorm.io/gorm"
)

type Login struct{}

type reqLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RespLogin struct {
	ArticleLikeSet []int64 `json:"articleLikeSet"`
	Avatar         string  `json:"avatar"`
	CommentLikeSet []int64 `json:"commentLikeSet"`
	Email          string  `json:"email"`
	Id             int64   `json:"id"`
	Intro          string  `json:"intro"`
	IpAddress      string  `json:"ipAddress"`
	IpSource       string  `json:"ipSource"`
	LastLoginTime  string  `json:"lastLoginTime"`
	LoginType      int     `json:"loginType"`
	Nickname       string  `json:"nickname"`
	TalkLikeSet    []int64 `json:"talkLikeSet"`
	UserInfoId     int64   `json:"userInfoId"`
	Username       string  `json:"username"`
}

type UserLikeRecord struct {
	ArticleLikeSet []int64 `json:"articleLikeSet"`
	CommentLikeSet []int64 `json:"commentLikeSet"`
	TalkLikeSet    []int64 `json:"talkLikeSet"`
}

func (l *Login) Login(ctx *gin.Context) {
	var form reqLogin
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "用户名和密码不能为空")
		return
	}
	curd := common.Curd{}
	user := common.TUserAuth{}
	ok, err := curd.Select(&user, "username = ?", form.Username)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	} else if !ok {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "不存在该用户名")
		return
	}
	// 验证密码
	if !common.VerifyPwd(user.Password, form.Password) {
		Response(ctx, errorcode.AuthorizedError, nil, false, "账号密码错误")
		return
	}
	// 通过IP获取地理位置
	ip := ctx.ClientIP()
	ipInfo, err := common.GetIpAddressAndSource(ip)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}

	//=====登陆成功
	userInfo := common.TUserInfo{}
	ok, err = curd.Select(&userInfo, "email = ?", user.Username)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	} else if !ok {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "不存在该用户信息")
		return
	}

	//更新用户最近登陆信息
	user.IpAddress = ip
	user.IpSource = ipInfo.Data.Province
	user.UserAgent = ctx.GetHeader("User-Agent")
	user.LastLoginTime = time.Now()
	user.LoginType = 1
	uaParse := user_agent.New(user.UserAgent)
	user.OS = fmt.Sprintf("%v", uaParse.OS())
	brows, _ := uaParse.Browser()
	user.Browser = fmt.Sprintf("%v", brows)
	ok, err = curd.Update(&common.TUserAuth{}, &user, "username = ?", user.Username)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	} else if !ok {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "不存在该用户")
		return
	}

	cache := &common.CacheOptions{
		Key:      fmt.Sprintf(rediskey.UserLike, userInfo.ID),
		Receiver: new(UserLikeRecord),
		Duration: rediskey.ExpireUserLike,
		Fun: func() (interface{}, error) {
			db := common.GetGorm()
			var likeRecord UserLikeRecord
			var tLike1 []common.TLike
			var tLike2 []common.TLike
			var tLike3 []common.TLike
			r1 := db.Model(&common.TLike{}).Where("user_id = ? and object = ?", userInfo.ID, "t_article").Find(&tLike1)
			r1 = db.Model(&common.TLike{}).Where("user_id = ? and object = ?", userInfo.ID, "t_talk").Find(&tLike2)
			r1 = db.Model(&common.TLike{}).Where("user_id = ? and object = ?", userInfo.ID, "t_comment").Find(&tLike3)
			if r1.Error != nil {
				return nil, err
			}
			for _, val := range tLike1 {
				likeRecord.ArticleLikeSet = append(likeRecord.ArticleLikeSet, val.ID)
			}
			for _, val := range tLike2 {
				likeRecord.TalkLikeSet = append(likeRecord.TalkLikeSet, val.ID)
			}
			for _, val := range tLike3 {
				likeRecord.CommentLikeSet = append(likeRecord.CommentLikeSet, val.ID)
			}

			return likeRecord, nil
		},
	}
	receiver, e1 := cache.GetSet()
	if e1 != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}

	data := RespLogin{}
	// 如果没有记录
	if receiver == nil {
		data.ArticleLikeSet = []int64{}
		data.CommentLikeSet = []int64{}
		data.TalkLikeSet = []int64{}
	} else {
		//如果有记录
		tl := receiver.(UserLikeRecord)
		data.ArticleLikeSet = tl.ArticleLikeSet
		data.CommentLikeSet = tl.CommentLikeSet
		data.TalkLikeSet = tl.TalkLikeSet
	}
	data.Id = user.ID
	data.Email = userInfo.Email
	data.Intro = userInfo.Intro
	data.IpAddress = user.IpAddress
	data.IpSource = user.IpSource
	data.Avatar = userInfo.Avatar
	data.LastLoginTime = user.LastLoginTime.String()
	data.LoginType = user.LoginType
	data.Nickname = userInfo.Nickname
	data.UserInfoId = user.UserInfoId
	data.Username = user.Username

	// 获取权限
	var userRole common.TUserRole
	db := common.GetGorm()
	r2 := db.Where("user_id = ?", userInfo.ID).First(&userRole)
	if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
		logger.Error(r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}

	token, err := auth.JwtEnc(user.ID, userRole.RoleId)
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "ticket",
		Value:  token,
		Secure: true,
	})
	if err != nil {
		err = errors.Wrap(err, "登陆生成token异常")
		logger.Error(fmt.Sprintf("%+v", err))
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 记录在线用户的auid
	redisClient := common.GetRedis()
	err = redisClient.SAdd(rediskey.OnlineUser, user.ID).Err()
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, &data, true, "操作成功")
	return

}

func (l *Login) LoginOut(ctx *gin.Context) {
	auid, _ := ctx.Get("a_userid")
	if auid != nil {
		// 下线用户
		redisClient := common.GetRedis()
		err := redisClient.SRem(rediskey.OnlineUser, auid.(int64)).Err()
		if err != nil && err != redis.Nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	// 删除cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "ticket",
		Value:  "",
		Secure: false,
		MaxAge: -1,
	})
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	ctrl "blog-server/controllers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Login struct {
	ctrl.LoginHandle
}

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

func (l *Login) Login(ctx *gin.Context) {
	var form reqLogin
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "用户名和密码不能为空")
		return
	}
	_session, err2 := Store.Get(ctx.Request, "CurUser")
	if err2 != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
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
	ok, err = curd.Select(&userInfo, "id = ?", user.ID)
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
		Receiver: new(common.TLike),
		Duration: rediskey.ExpireUserLike,
		Fun: func() (interface{}, error) {
			var tLike = &common.TLike{}
			ok, err = curd.Select(tLike, "user_id = ?", user.ID)
			if err != nil {
				return nil, err
			} else if !ok {
				return nil, nil
			}
			return tLike, nil
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
		tl := struct {
			ArticleLikeSet []int64 `json:"articleLikeSet"`
			CommentLikeSet []int64 `json:"commentLikeSet"`
			TalkLikeSet    []int64 `json:"talkLikeSet"`
		}{}
		r := receiver.(*common.TLike)
		err = json.Unmarshal([]byte(r.LikeItem), &tl)
		if err != nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		data.ArticleLikeSet = tl.ArticleLikeSet
		data.CommentLikeSet = tl.CommentLikeSet
		data.TalkLikeSet = tl.TalkLikeSet
	}
	data.Id = user.ID
	data.Email = userInfo.Email
	data.Intro = userInfo.Intro
	data.IpAddress = user.IpAddress
	data.IpSource = user.IpSource
	data.LastLoginTime = user.LastLoginTime.String()
	data.LoginType = user.LoginType
	data.Nickname = userInfo.Nickname
	data.UserInfoId = user.UserInfoId
	data.Username = user.Username

	_session.Values["userid"] = user.ID
	_session.Values["login_time"] = time.Now().Unix()
	err = _session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, &data, true, "操作成功")
	return

}

func (l *Login) LoginOut(ctx *gin.Context) {

}

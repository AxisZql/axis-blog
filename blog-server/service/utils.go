package service

import (
	"blog-server/common/errorcode"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
* @author:AxisZql
* @date: 2022-3-16 11:33 PM
* @desc: 响应统一处理模块
 */

func Response(ctx *gin.Context, code int64, data interface{}, flag bool, message string) {
	resp := gin.H{
		"code":    code,
		"data":    &data,
		"flag":    flag,
		"message": message,
	}
	ctx.JSON(http.StatusOK, resp)
}

/*
* @author:AxisZql
* @date:2022-3-15 3:32 PM
* @desc:应用中间件
 */

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// before request
		_session, err := Store.Get(ctx.Request, "CurUser")
		if err != nil {
			// 对下一步处理函数对执行进行拦截
			ctx.Abort()
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		userid := _session.Values["a_userid"]
		b := _session.Values["login_time"]
		if userid == nil || b == nil {
			ctx.Abort()
			delete(_session.Values, "CurUser")
			_ = _session.Save(ctx.Request, ctx.Writer)
			Response(ctx, errorcode.AuthorizedError, nil, false, "没有操作权限")
			return
		}
		a := time.Now().Unix()
		_b := b.(int64)
		if a-_b >= 30*60 {
			ctx.Abort()
			delete(_session.Values, "CurUser")
			_ = _session.Save(ctx.Request, ctx.Writer)
			Response(ctx, errorcode.ExpireLoginTime, nil, false, "登陆状态过期")
			return
		}

		ctx.Next()
		// response
	}
}

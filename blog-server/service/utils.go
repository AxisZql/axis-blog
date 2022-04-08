package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
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
		if a-_b >= 7*24*60 {
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

// == 请求日志中间件
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func LogAopReq() func(c *gin.Context) {
	return func(c *gin.Context) {
		db := common.GetGorm()
		_session, _ := Store.Get(c.Request, "CurUser")
		userid := _session.Values["a_userid"]
		var ua common.TUserAuth
		var ui common.TUserInfo
		r1 := db.Where("id = ?", userid).First(&ua)
		r1 = db.Where("id = ?", ua.UserInfoId).First(&ui)
		if r1.Error != nil {
			c.Abort()
			logger.Error(r1.Error.Error())
			Response(c, errorcode.AuthorizedError, nil, false, "没有操作权限")
			return
		}
		// 请求参数
		req, err := c.GetRawData()
		if err != nil {
			c.Abort()
			logger.Error(err.Error())
			Response(c, errorcode.Fail, nil, false, "系统异常")
			return
		}
		reqParam := fmt.Sprintf("%v", string(req))
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req)) //由于body只能读一次，为了防止处理函数不能正确获取body故不关闭body

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		// 处理请求
		c.Next()
		if strings.Contains(c.Request.URL.Path, "images") || strings.Contains(c.Request.URL.Path, "cover") || strings.Contains(c.Request.URL.Path, "logs") {
			return
		}
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 请求IP
		clientIP := c.ClientIP()
		ipInfo, _ := common.GetIpAddressAndSource(clientIP)
		path := c.Request.URL.Path
		pL := strings.Split(path, "/")
		h := c.HandlerNames()
		funcName := h[len(h)-1]
		opName := strings.Split(funcName, ".")
		l := common.TOperationLog{
			OptModule:     pL[len(pL)-1],
			OptType:       reqMethod,
			OptUrl:        reqUri,
			OptMethod:     funcName,
			OptDesc:       opName[len(opName)-1],
			RequestParam:  reqParam,
			RequestMethod: reqMethod,
			ResponseData:  w.body.String(),
			UserId:        ui.ID,
			Nickname:      ui.Nickname,
			IpAddress:     clientIP,
			IpSource:      ipInfo.Data.Province,
		}
		r1 = db.Model(&common.TOperationLog{}).Create(&l)
		if r1.Error != nil {
			c.Abort()
			logger.Error(r1.Error.Error())
			Response(c, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
}

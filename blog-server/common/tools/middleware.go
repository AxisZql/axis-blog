package tools

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date:2022-3-15 3:32 PM
* @desc:应用中间件
 */

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// before request
		ctx.Next()
		// response
	}
}

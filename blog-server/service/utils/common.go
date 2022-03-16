package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

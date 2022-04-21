package controllers

import "github.com/gin-gonic/gin"

/*
@author AxisZql
@desc websocket 实现实时在线聊天
@date 2022-4-21 3.32 PM
*/

type WebSocket interface {
	WebSocketHandle(ctx *gin.Context)
}

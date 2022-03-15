package controllers

import "github.com/gin-gonic/gin"

/*
* @author:AxisZql
* @date: 2022-3-15 7:22 PM
* @desc: 说说模块接口
 */

type TalkHandle interface {
	ListHomeTalks(*gin.Context)    //查看首页说说
	ListTalks(*gin.Context)        //查看说说列表
	GetTalkById(*gin.Context)      //根据id查看说说
	SaveTalkLike(*gin.Context)     //点赞说说
	SaveTalkImages(*gin.Context)   //上传说说图片
	SaveOrUpdateTalk(*gin.Context) //保存或者修改说说
	DeleteTalks(*gin.Context)      //删除说说
	ListBackTalks(*gin.Context)    //查看后台说说
	GetBackTalkById(*gin.Context)  //根据id查看后台说说
}

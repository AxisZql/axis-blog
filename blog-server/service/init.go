package service

import (
	"blog-server/common"
	"blog-server/common/tools"

	"github.com/gorilla/sessions"
)

var logger = common.GetLogger()
var Store = sessions.NewCookieStore([]byte(common.Conf.Jwt.Key))
var senitiveForest = tools.GetSenitiveForest()

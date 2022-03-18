package service

import (
	"blog-server/common"
	"github.com/gorilla/sessions"
)

var logger = common.GetLogger()
var Store = sessions.NewCookieStore([]byte(common.Conf.Jwt.Key))

package common

import (
	"fmt"
	"os"
)

var logger = Logger()

func EnvInit() {
	//初始化配置文件
	err := InitConfigure()
	if err != nil {
		logger.Error(fmt.Sprintf("初始化配置文件失败:%v", err))
	}

	err = os.MkdirAll(Conf.App.ConfigDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建配置图片目录失败:%v", err))
	}
	err = os.MkdirAll(Conf.App.AvatarDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建头像目录失败:%v", err))
	}
	err = os.MkdirAll(Conf.App.ArticleDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建文章图片目录失败:%v", err))
	}
	err = os.MkdirAll(Conf.App.VoiceDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建语音目录失败:%v", err))
	}
	err = os.MkdirAll(Conf.App.PhotoDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建相册目录失败:%v", err))
	}
	err = os.MkdirAll(Conf.App.TalkDir, os.ModePerm)
	if err != nil {
		logger.Error(fmt.Sprintf("创建说说图片目录失败:%v", err))
	}

	// 初始化数据库
	if err = InitDb(); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("=========== finish to init app env ===========")

}

package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// Logger 自定义日志记录器
func Logger() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   "/tmp/axis-blog/log.log", //日志文件路径
		MaxSize:    128,                      //每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                       //日志文件最多保存多少个备份
		MaxAge:     7,                        //文件最多保存多少天
		Compress:   true,                     //是否压缩
	}

	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    //ISO8601 UTC时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel) //设置日志级别

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encodeConfig),                                            // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), //打印日志到控制台和文件
		atomicLevel, //日志级别
	)
	//开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	logger := zap.New(core, caller, development)
	return logger
}

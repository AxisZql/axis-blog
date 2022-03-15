package main

import (
	"blog-server/cmd"
	"fmt"
	"runtime"
)

var (
	_version   = "unknown"
	_gitCommit = "unknown"
	_goVersion = runtime.Version() //go 版本
	_buildTime = "unknown"
	_osArch    = runtime.GOARCH //系统架构
)

var description string

// 描述信息初始化
func init() {
	description = fmt.Sprintf(`BUILD INFO::
    Version:		%s
    Go Version:  	%s
    Git Commit:  	%s
    Build time: 	%s
    OS/Arch:   		%s`, _version, _goVersion, _gitCommit, _buildTime, _osArch)
}

func main() {
	fmt.Println(description)
	cmd.Execute()
}

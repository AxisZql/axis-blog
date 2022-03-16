package common

/*
* @author:AxisZql
* @date: 2022-3-16 4:21 PM
* @desc: 查看数据库接口封装模块
 */

type CurdHandle interface {
	Select(table interface{}, condition string, val ...interface{}) (bool, error)
	Update(table interface{}, condition string, val ...interface{}) (bool, error)
	Create(table interface{}, field ...string) error
}

package common

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Curd struct {
	CurdHandle
}

/*
	@table 表结构
	@val 查询值
	@query 查询目标属性
*/
func (c *Curd) Select(table interface{}, condition string, val ...interface{}) (bool, error) {
	db := GetGorm()
	result := db.Where(condition, val).First(table)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Sprintf("查询数据库失败:%s", result.Error))
		return false, result.Error
	}
	if result.Error == nil {
		return true, nil
	}
	return false, nil

}

func (c *Curd) Update(table interface{}, val interface{}, condition string, queryVal ...interface{}) (bool, error) {
	db := GetGorm()
	result := db.Model(table).Where(condition, queryVal).Updates(val)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Sprintf("更新数据库失败:%s", result.Error))
		return false, result.Error
	}
	if result.Error == nil {
		return true, nil
	}
	return false, nil
}

func (c *Curd) Create(table interface{}, field ...string) error {
	db := GetGorm()
	result := db.Select(field).Create(table)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *Curd) SqlQuery(sql string, dest interface{}, val ...interface{}) (bool, error) {
	db := GetGorm()
	result := db.Exec(sql, val).First(&dest)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Sprintf("SQL语句执行失败:%v", result.Error))
		return false, result.Error
	}
	if result.Error == nil {
		return true, nil
	}
	return false, nil
}

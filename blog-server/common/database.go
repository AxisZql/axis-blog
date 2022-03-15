package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

/*
* @author:AxisZql
* @date:2022-3-15 1:59 PM
* @desc:数据库初始化
 */

var (
	db *gorm.DB
)

func InitDb() error {
	d := Conf.Db
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.DbName)

	var err error
	db, err = gorm.Open(mysql.Open(dbInfo), &gorm.Config{
		Logger: glog.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			glog.Config{
				SlowThreshold: time.Second,
				LogLevel:      glog.Info,
				Colorful:      true,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名
		},
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("数据库:%s:%s:%s;%s", d.Type, d.Host, d.Port, d.DbName))
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	//设置连接池中最大大闲连接数
	sqlDB.SetMaxIdleConns(10)
	//设置数据库大最大连接数
	sqlDB.SetMaxOpenConns(5)
	//设置连接大最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	if Conf.App.InitModels {
		t := time.Now()
		modelsInit()
		logger.Info(fmt.Sprintf("inti models in:%v", time.Since(t)))
	}
	return nil
}

func modelsInit() {
	logger.Info("models initializing...")
	t := time.Now()
	e1 := db.AutoMigrate(&TArticle{}, &TCategory{}, &TTag{}, &TArticleTag{}, &TChatRecord{}, &TComment{}, &TFriendLink{},
		&TMenu{}, &TMessage{}, &TOperationLog{}, &TPage{}, &TPhoto{}, &TPhotoAlbum{}, &TResource{}, &TRole{}, &TRoleMenu{},
		&TRoleResource{}, &TTalk{}, &TUniqueView{}, &TUserAuth{}, &TUserInfo{}, &TUserRole{}, &TWebsiteConfig{})
	if e1 != nil {
		err := fmt.Errorf("初始化表失败:%v", e1)
		panic(err)
	}
	logger.Debug(fmt.Sprintf("models inited in:%s", time.Since(t)))
}

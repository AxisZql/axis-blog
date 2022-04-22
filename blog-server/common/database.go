package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

/*
* @author:AxisZql
* @date:2022-3-15 1:59 PM
* @desc:数据库初始化
 */

var (
	// user后台菜单组件权限表
	vUserMenu = `create view v_user_menu as
       select user_id,
       menu_id,
       name,
       path,
       component,
       icon,
       order_num,
       parent_id,
       is_hidden
from t_user_role ur
         inner join t_role_menu rm
                    on ur.role_id = rm.role_id
         inner join t_menu m on rm.menu_id = m.id order by menu_id asc;`
	// 类别文章数量视图
	vCategoryCount = `
     create view v_category_count as
     select tc.id, tc.category_name, count(*) as count
     from t_category tc,
     t_article ta
     where tc.id = ta.category_id
     GROUP by tc.category_name, tc.id;`
	// 统计每一天的文章数量
	vArticleStatistics = `
    CREATE VIEW v_article_statistics as
    SELECT DATE(create_time) AS date, COUNT(*) AS count
    FROM t_article
    GROUP BY date
    ORDER BY date desc ;`
	// 文章信息视图
	vArticleInfo = `
    create view v_article_info as
SELECT a.id,
       article_cover,
       article_title,
       article_content,
       a.type,
       a.is_top,
       a.category_id,
       category_name,
       view_count,
       like_count,
       original_url,
       a.create_time,
       a.update_time
FROM (SELECT t_article.id,
             article_cover,
             article_title,
             article_content,
             type,
             is_top,
             view_count,
             tc.like_count,
             original_url,
             create_time,
             update_time,
             category_id
      FROM t_article
               left join (select COUNT(*) as like_count, t_article.id
                          from t_like,
                               t_article
                          where t_like.like_id = t_article.id
                            and t_like.object = 't_article'
                          GROUP BY t_article.id) as tc on tc.id = t_article.id
      WHERE is_delete = 0
        AND status = 1
      ORDER BY is_top DESC
     ) a
         inner JOIN t_category c ON a.category_id = c.id
ORDER BY a.is_top DESC, a.update_time DESC, a.create_time DESC;`
	// 评论视图
	vComment = `
    create view v_comment as 
    select ans.*, tu2.nickname as reply_nickname, tu2.web_site as reply_web_site
from (select tc.id,
             parent_id,
             tu.id as user_id,
             nickname,
             avatar,
             web_site,
             reply_user_id,
             comment_content,
             topic_id,
             type,
             is_delete,
             is_review,
             tc.create_time
      from (select id, avatar, nickname, web_site
            from t_user_info) tu
               join t_comment tc on tu.id = tc.user_id) ans
         left join t_user_info tu2
                   on tu2.id = ans.reply_user_id
where is_review = 1
  and is_delete = 0;`
	// 说说视图
	vTalkInfo = `
create view v_talk_info as 
select t1.*, lc.like_count, lco.comment_count
from (select  tt.*,avatar,nickname
      from t_user_info tu
               join t_talk tt on tu.id = tt.user_id)
         as t1
         left join

     (select count(*) as like_count, t_talk.id
      from t_talk,
           t_like
      where t_like.like_id = t_talk.id
        and t_like.object = 't_talk' group by t_talk.id) as lc on t1.id = lc.id
         left join
     (select count(*) as comment_count, t_talk.id
      from t_comment,
           t_talk
      where t_comment.topic_id = t_talk.id
        and t_comment.type = 3 group by t_talk.id) lco on t1.id = lco.id order by is_top DESC;
`
	vUserInfo = `
		create view v_user_info as 
select tua.*, tui.email, tui.nickname, tui.avatar, tui.intro, tui.web_site, tui.is_disable
from (select id,
             user_info_id,
             username,
             login_type,
             last_login_time,
             ip_source,
             ip_address,
             user_agent,
             os,
             browser,
             create_time
      from t_user_auth) tua
         join
     (select id, email, nickname, avatar, intro, web_site, is_disable
      from t_user_info) tui on tua.user_info_id = tui.id
`
)

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
	e1 := db.AutoMigrate(&TArticle{}, &TLike{}, &TCategory{}, &TTag{}, &TArticleTag{}, &TChatRecord{}, &TComment{}, &TFriendLink{},
		&TMenu{}, &TMessage{}, &TOperationLog{}, &TPage{}, &TPhoto{}, &TPhotoAlbum{}, &TResource{}, &TRole{}, &TRoleMenu{},
		&TRoleResource{}, &TTalk{}, &TUniqueView{}, &TUserAuth{}, &TUserInfo{}, &TUserRole{}, &TWebsiteConfig{})
	if e1 != nil {
		err := fmt.Errorf("初始化表失败:%v", e1)
		panic(err)
	}
	e2 := db.Exec("drop view v_user_menu;")
	if e2.Error != nil && !strings.Contains(e2.Error.Error(), "Unknown table") {
		panic(e2)
	}
	e2 = db.Exec(vUserMenu)
	if e2.Error != nil {
		panic(e2)
	}
	e3 := db.Exec("drop view v_category_count;")
	if e3.Error != nil && !strings.Contains(e3.Error.Error(), "Unknown table") {
		panic(e2)
	}
	e3 = db.Exec(vCategoryCount)
	if e3.Error != nil {
		panic(e3)
	}
	e4 := db.Exec("drop view v_article_statistics;")
	if e4.Error != nil && !strings.Contains(e4.Error.Error(), "Unknown table") {
		panic(e2)
	}
	e4 = db.Exec(vArticleStatistics)
	if e4.Error != nil {
		panic(e4)
	}
	e6 := db.Exec("drop view v_article_info;")
	if e6.Error != nil && !strings.Contains(e6.Error.Error(), "Unknown table") {
		panic(e2)
	}
	e6 = db.Exec(vArticleInfo)
	if e6.Error != nil {
		panic(e6)
	}
	e7 := db.Exec("drop view v_comment;")
	if e7.Error != nil && !strings.Contains(e7.Error.Error(), "Unknown table") {
		panic(e2)
	}
	e7 = db.Exec(vComment)
	if e7.Error != nil {
		panic(e7)
	}
	e8 := db.Exec("drop view v_talk_info")
	if e8.Error != nil && !strings.Contains(e8.Error.Error(), "Unknown table") {
		panic(e8)
	}
	e8 = db.Exec(vTalkInfo)
	if e8.Error != nil {
		panic(e8)
	}
	e9 := db.Exec("drop view v_user_info")
	if e9.Error != nil && !strings.Contains(e9.Error.Error(), "Unknown table") {
		panic(e9)
	}
	e9 = db.Exec(vUserInfo)
	if e9.Error != nil {
		panic(e9)
	}
	logger.Debug(fmt.Sprintf("models inited in:%s", time.Since(t)))
}

func GetGorm() *gorm.DB {
	return db
}

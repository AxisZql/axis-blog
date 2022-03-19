package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	ctrl "blog-server/controllers"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BlogInfo struct {
	ctrl.BlogInfo
}

func (b *BlogInfo) GetBlogHomeInfo(*gin.Context) {}

type CategoryDto struct {
	Id           int64  `json:"id"`
	CategoryName string `json:"categoryName"`
	Count        int64  `json:"articleCount"`
}

type TagDto struct {
	Id      int64  `json:"id"`
	TagName string `json:"tagName"`
}

type ArticleStatistics struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
type UniqueViewsDTO struct {
	Day        time.Time `json:"day"`
	ViewsCount int64     `json:"viewsCount"`
}
type ArticleRankDTO struct {
	ArticleTitle string `json:"articleTitle"`
	ViewCount    int64  `json:"viewsCount"`
}

type respGetBlogBackInfo struct {
	ViewsCount            int64               `json:"viewsCount"`
	MessageCount          int64               `json:"messageCount"`
	UserCount             int64               `json:"userCount"`
	ArticleCount          int64               `json:"articleCount"`
	CategoryDTOList       []CategoryDto       `json:"categoryDTOList"`
	TagDTOList            []TagDto            `json:"tagDTOList"`
	ArticleStatisticsList []ArticleStatistics `json:"articleStatisticsList"`
	UniqueViewDTOList     []UniqueViewsDTO    `json:"uniqueViewDTOList"`
	ArticleRankDTOList    []ArticleRankDTO    `json:"articleRankDTOList"`
}

func (b *BlogInfo) GetBlogBackInfo(ctx *gin.Context) {
	data := respGetBlogBackInfo{}
	db := common.GetGorm()
	//查询博客访问量
	rc := common.GetRedis()
	viewCount, err := rc.Get(rediskey.BlogViewsCount).Int64()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data.ViewsCount = viewCount
	//获取留言数量
	var messageCount int64
	// 获取用户数量
	var userCount int64
	// 获取文章数量
	var articleCount int64
	// 获取标签列表
	var tagList []TagDto
	r1 := db.Model(&common.TMessage{}).Count(&messageCount)
	r2 := db.Model(&common.TUserAuth{}).Count(&userCount)
	r3 := db.Model(&common.TArticle{}).Count(&articleCount)
	r4 := db.Find(&[]common.TTag{}).Scan(&tagList)
	if r1.Error != nil || r2.Error != nil || r3.Error != nil || r4.Error != nil {
		logger.Error(r1.Error.Error() + "|||" + r2.Error.Error() + "|||" + r3.Error.Error() + "|||" + r4.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data.MessageCount = messageCount
	data.UserCount = userCount
	data.ArticleCount = articleCount
	data.TagDTOList = tagList
	// 获取文章分类列表
	categoryDTOList := &data.CategoryDTOList
	rows, err := db.Raw("select * from v_category_count").Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for rows.Next() {
		var t CategoryDto
		_ = db.ScanRows(rows, &t)
		*categoryDTOList = append(*categoryDTOList, t)
	}
	// 获取文章分布日期统计列表
	articleStatistics := &data.ArticleStatisticsList
	rows, err = db.Raw("select * from v_article_statistics;").Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for rows.Next() {
		var t ArticleStatistics
		_ = db.ScanRows(rows, &t)
		*articleStatistics = append(*articleStatistics, t)
	}
	// 获取最近一周的访客数量
	currentTime := time.Now()
	lastWeek := currentTime.AddDate(0, 0, -7) //获取一周前的时间
	var count int64
	r5 := db.Raw("select SUM(views_count) as count from t_unique_view where create_time > ? and create_time <= ?", lastWeek, currentTime).Scan(&count)
	if r5.Error != nil && r5.Error != gorm.ErrRecordNotFound {
		logger.Error(r5.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data.UniqueViewDTOList = append(data.UniqueViewDTOList, UniqueViewsDTO{ViewsCount: count, Day: lastWeek})
	//查看数据库中访问量前5的文章
	articleRankDTOList := &data.ArticleRankDTOList
	rows, err = db.Raw("select * from t_article order by view_count desc limit 5;").Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for rows.Next() {
		var t ArticleRankDTO
		_ = db.ScanRows(rows, &t)
		*articleRankDTOList = append(*articleRankDTOList, t)
	}
	Response(ctx, errorcode.Success, data, true, "操作成功")

}
func (b *BlogInfo) SavePhotoAlbumCover(*gin.Context) {}
func (b *BlogInfo) UpdateWebsiteConfig(*gin.Context) {}
func (b *BlogInfo) GetWebSiteConfig(*gin.Context)    {}
func (b *BlogInfo) GetAbout(*gin.Context)            {}
func (b *BlogInfo) UpdateAbout(*gin.Context)         {}
func (b *BlogInfo) SendVoice(*gin.Context)           {}

// Report 报告唯一访客信息
func (b *BlogInfo) Report(ctx *gin.Context) {
	ip := ctx.ClientIP()
	ipSource, err := common.GetIpAddressAndSource(ip)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	payload := ip + userAgent

	md := fmt.Sprintf("%x", md5.Sum([]byte(payload)))
	// 获取用户唯一标识，如不存在则新建
	exist, err := common.GetSetSetCache(rediskey.UniqueVisitor, md)
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 如果是新访客
	if !exist {
		redisClient := common.GetRedis()
		result, err := redisClient.HGet(rediskey.VisitorArea, ipSource.Data.Province).Result()
		if err != nil && err != redis.Nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if err == redis.Nil {
			err := redisClient.HSet(rediskey.VisitorArea, ipSource.Data.Province, 1).Err()
			if err != nil {
				logger.Error(err.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		} else {
			b, _ := strconv.Atoi(result)
			err = redisClient.HSet(rediskey.VisitorArea, ipSource.Data.Province, b).Err()
			if err != nil {
				logger.Error(err.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		// 更新博客浏览人数，先更新数据库中的数据
		curd := common.Curd{}
		uView := common.TUniqueView{}
		exist, err := curd.Select(&uView, "to_days(create_time) = to_days(?)", time.Now())
		if err != nil {
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if exist {
			_, e := curd.Update(&common.TUniqueView{}, common.TUniqueView{ViewsCount: uView.ViewsCount + 1, UpdateTime: time.Now()}, "to_days(create_time) = to_days(?)", uView.CreateTime)
			if e != nil {
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		} else {
			uView.ViewsCount = 1
			e := curd.Create(&uView, "views_count")
			if e != nil {
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		err = redisClient.Incr(rediskey.BlogViewsCount).Err()
		if err != nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

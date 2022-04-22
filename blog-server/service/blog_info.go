package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"blog-server/common/rediskey"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"time"
)

type BlogInfo struct{}

type page struct {
	PageCover string `json:"pageCover"`
	ID        int64  `json:"id"`
	PageName  string `json:"pageName"`
	PageLabel string `json:"pageLabel"`
}

type webConfig struct {
	WebsiteAvatar     string   `json:"websiteAvatar"`
	WebsiteName       string   `json:"websiteName"`
	WebsiteAuthor     string   `json:"websiteAuthor"`
	WebsiteIntro      string   `json:"websiteIntro"`
	WebsiteNotice     string   `json:"websiteNotice"`
	WebsiteCreateTime string   `json:"websiteCreateTime"`
	WebsiteRecordNo   string   `json:"websiteRecordNo"`
	SocialLoginList   []string `json:"socialLoginList"`
	SocialUrlList     []string `json:"socialUrlList"`
	QQ                string   `json:"qq"`
	Github            string   `json:"github"`
	Gitee             string   `json:"gitee"`
	TouristAvatar     string   `json:"touristAvatar"`
	UserAvatar        string   `json:"userAvatar"`
	IsCommentReview   int      `json:"isCommentReview"`
	IsMessageReview   int      `json:"isMessageReview"`
	WebsocketUrl      string   `json:"websocketUrl"`
	IsEmailNotice     int      `json:"isEmailNotice"`
	IsReward          int      `json:"isReward"`
	WeiXinQRCode      string   `json:"weiXinQRCode"`
	AlipayQRCode      string   `json:"alipayQRCode"`
	IsChatRoom        int      `json:"isChatRoom"`
	IsMusicPlayer     int      `json:"isMusicPlayer"`
}

type respBlogHomeInfo struct {
	ArticleCount  int64     `json:"articleCount"`
	CategoryCount int64     `json:"categoryCount"`
	TagCount      int64     `json:"tagCount"`
	ViewsCount    int64     `json:"viewsCount"`
	WebsiteConfig webConfig `json:"websiteConfig"`
	PageList      []page    `json:"pageList"`
}

func (b *BlogInfo) GetBlogHomeInfo(ctx *gin.Context) {
	//文章数量
	var articleCount int64
	// 分类数量
	var categoryCount int64
	// 标签数量
	var tagCount int64
	db := common.GetGorm()
	r1 := db.Model(&common.TArticle{}).Count(&articleCount)
	r2 := db.Model(&common.TCategory{}).Count(&categoryCount)
	r3 := db.Model(&common.TTag{}).Count(&tagCount)
	if r1.Error != nil || r2.Error != nil || r3.Error != nil {
		logger.Error(r1.Error.Error() + "|||" + r2.Error.Error() + "|||" + r3.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	redisClient := common.GetRedis()
	viewsCount, err := redisClient.Get(rediskey.BlogViewsCount).Int64()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	webConfigStr, e := redisClient.Get(rediskey.WebsiteConfig).Result()
	if e != nil && e != redis.Nil {
		logger.Error(e.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if e == redis.Nil {
		// 获取网站配置
		var wConfig common.TWebsiteConfig
		r4 := db.Model(&common.TWebsiteConfig{}).First(&wConfig)
		if r4.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		webConfigStr = wConfig.Config
		// 将网站配置保存到redis中
		if err2 := redisClient.Set(rediskey.WebsiteConfig, wConfig.Config, -1).Err(); err2 != nil {
			logger.Error(err2.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	pageList := make([]page, 0)
	r5 := db.Model(&common.TPage{}).Find(&pageList)
	if r5.Error != nil {
		logger.Error(r5.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := respBlogHomeInfo{}
	data.ArticleCount = articleCount
	data.CategoryCount = categoryCount
	data.TagCount = tagCount
	data.ViewsCount = viewsCount
	w := webConfig{}
	_ = json.Unmarshal([]byte(webConfigStr), &w)
	data.WebsiteConfig = w
	data.PageList = pageList

	Response(ctx, errorcode.Success, data, true, "操作成功")

}

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

type reqSaveConfigPic struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (b *BlogInfo) SaveConfigPic(ctx *gin.Context) {
	var form reqSaveConfigPic
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	f, _ := form.File.Open()
	extendName := strings.Split(form.File.Filename, ".")
	if len(extendName) != 2 && extendName[1] != "png" && extendName[1] != "gif" && extendName[1] != "jpg" {
		Response(ctx, errorcode.ValidError, nil, false, "不支持的图片格式;仅支持png|gif|jpg格式")
		return
	}
	defer f.Close()
	fileData, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		logger.Error(err2.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	fileMD5 := fmt.Sprintf("%x", md5.Sum(fileData))
	fileName := fileMD5 + "." + extendName[1]
	filePath := common.Conf.App.ConfigDir + fileName
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	imgUrl := fmt.Sprintf("%s/config/%s", common.Conf.App.HostName, fileName)
	Response(ctx, errorcode.Fail, imgUrl, true, "操作成功")
}
func (b *BlogInfo) UpdateWebsiteConfig(ctx *gin.Context) {
	var form webConfig
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var wf common.TWebsiteConfig
	r1 := db.Model(&common.TWebsiteConfig{}).First(&wf)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data, err := json.Marshal(form)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	wf.Config = string(data)
	wf.UpdateTime = time.Now()
	r1 = db.Save(&wf)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	redisClient := common.GetRedis()
	e := redisClient.Set(rediskey.WebsiteConfig, data, -1).Err()
	if e != nil {
		logger.Error(e.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (b *BlogInfo) GetWebSiteConfig(ctx *gin.Context) {
	redisClient := common.GetRedis()
	w, err := redisClient.Get(rediskey.WebsiteConfig).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if err == redis.Nil {
		db := common.GetGorm()
		var wf common.TWebsiteConfig
		r1 := db.Model(&common.TWebsiteConfig{}).First(&wf)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		e := redisClient.Set(rediskey.WebsiteConfig, wf.Config, -1).Err()
		if e != nil {
			logger.Error(e.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		w = wf.Config
	}
	var _w webConfig
	err = json.Unmarshal([]byte(w), &_w)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, _w, true, "操作成功")
}
func (b *BlogInfo) GetAbout(ctx *gin.Context) {
	redisClient := common.GetRedis()
	about, err := redisClient.Get(rediskey.ABOUT).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, about, true, "操作成功")
}

type reqUpdateAbout struct {
	AboutContent string `json:"aboutContent" binding:"required"`
}

func (b *BlogInfo) UpdateAbout(ctx *gin.Context) {
	var form reqUpdateAbout
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	redisClient := common.GetRedis()
	err := redisClient.Set(rediskey.ABOUT, form.AboutContent, -1).Err()
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqSendVoice struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Type      int                   `form:"type" binding:"required"`
	Nickname  string                `form:"nickname" `
	Avatar    string                `form:"avatar"`
	IpAddress string                `form:"ipAddress"`
	IpSource  string                `form:"ipSource"`
}

func (b *BlogInfo) SendVoice(ctx *gin.Context) {
	var form reqSendVoice
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	f, _ := form.File.Open()
	extendName := strings.Split(form.File.Filename, ".")
	if len(extendName) != 2 && extendName[1] != "wav" {
		Response(ctx, errorcode.ValidError, nil, false, "不支持的语音格式;仅支持wav格式")
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(f)
	fileData, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		logger.Error(err2.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	fileMD5 := fmt.Sprintf("%x", md5.Sum(fileData))
	fileName := fileMD5 + "." + extendName[1]
	filePath := common.Conf.App.VoiceDir + fileName
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	voiceUrl := fmt.Sprintf("%s/voice/%s", common.Conf.App.HostName, fileName)
	Response(ctx, errorcode.Success, nil, true, "操作成功")
	db := common.GetGorm()
	chat := common.TChatRecord{
		Type:      form.Type,
		Avatar:    form.Avatar,
		Content:   voiceUrl,
		IpAddress: form.IpAddress,
		IpSource:  form.IpSource,
		Nickname:  form.Nickname,
	}
	r1 := db.Model(&common.TChatRecord{}).Create(&chat)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
	}
	if r1.Error == nil {
		resp := WsMessage{
			Type: form.Type,
			Data: chat,
		}
		msg, _ := json.Marshal(&resp)
		// 把语音消息群发回去
		Manager.Broadcast <- msg
	}
}

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
		result, err := redisClient.HGet(rediskey.VisitorArea, ipSource.Data.Province).Int64()
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
			b := result + 1
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

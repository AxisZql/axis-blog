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
	"strconv"
)

type BlogInfo struct {
	ctrl.BlogInfo
}

func (b *BlogInfo) GetBlogHomeInfo(*gin.Context)     {}
func (b *BlogInfo) GetBlogBackInfo(*gin.Context)     {}
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
		// 更新博客浏览人数
		err = redisClient.Incr(rediskey.BlogViewsCount).Err()
		if err != nil {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

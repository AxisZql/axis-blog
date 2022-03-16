package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
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
func (b *BlogInfo) Report(*gin.Context)              {}

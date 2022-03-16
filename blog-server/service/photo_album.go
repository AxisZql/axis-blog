package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type PhotoAlbum struct {
	ctrl.PhotoAlbumHandle
}

func (p *PhotoAlbum) SavePhotoAlbumCover(*gin.Context)    {}
func (p *PhotoAlbum) SaveOrUpdatePhotoAlbum(*gin.Context) {}
func (p *PhotoAlbum) ListPhotoAlbumBack(*gin.Context)     {}
func (p *PhotoAlbum) ListPhotoAlbumBackInfo(*gin.Context) {}
func (p *PhotoAlbum) GetPhotoAlbumBackById(*gin.Context)  {}
func (p *PhotoAlbum) DeletePhotoAlbumById(*gin.Context)   {}
func (p *PhotoAlbum) ListPhotoAlbum(*gin.Context)         {}

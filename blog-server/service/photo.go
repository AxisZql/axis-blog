package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Photo struct {
	ctrl.PhotoHandle
}

func (p *Photo) ListPhotos(*gin.Context)         {}
func (p *Photo) UpdatePhoto(*gin.Context)        {}
func (p *Photo) SavePhoto(*gin.Context)          {}
func (p *Photo) UpdatePhotoAlbum(*gin.Context)   {}
func (p *Photo) UpdatePhotoDelete(*gin.Context)  {}
func (p *Photo) DeletePhotos(*gin.Context)       {}
func (p *Photo) ListPhotoByAlbumId(*gin.Context) {}

package service

import (
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Article struct {
	ctrl.ArticleHandle
}

func (a *Article) ListArchives(ctx *gin.Context)           {}
func (a *Article) ListArticles(ctx *gin.Context)           {}
func (a *Article) ListArticleBack(ctx *gin.Context)        {}
func (a *Article) SaveOrUpdateArticle(ctx *gin.Context)    {}
func (a *Article) UpdateArticleTop(ctx *gin.Context)       {}
func (a *Article) UpdateArticleDelete(ctx *gin.Context)    {}
func (a *Article) SaveArticleImages(ctx *gin.Context)      {}
func (a *Article) DeleteArticle(ctx *gin.Context)          {}
func (a *Article) GetArticleBackById(ctx *gin.Context)     {}
func (a *Article) GetArticleById(ctx *gin.Context)         {}
func (a *Article) ListArticleByCondition(ctx *gin.Context) {}
func (a *Article) ListArticleBySearch(ctx *gin.Context)    {}
func (a *Article) SaveArticleLike(ctx *gin.Context)        {}

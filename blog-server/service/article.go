package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Article struct{}

type reqListArchives struct {
	Current int `form:"current" binding:"required"`
}
type archive struct {
	ID           int64     `json:"id"`
	ArticleTitle string    `json:"articleTitle"`
	CreateTime   time.Time `json:"createTime"`
}

func (a *Article) ListArchives(ctx *gin.Context) {
	var form reqListArchives
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 {
		form.Current = 1
	}
	db := common.GetGorm()
	var articleCount int64
	var archiveList []archive
	r1 := db.Model(&common.TArticle{}).Count(&articleCount)
	r2 := db.Model(&common.TArticle{}).Limit(10).Offset((form.Current - 1) * 10).Order("create_time DESC").Find(&archiveList)
	if r1.Error != nil || r2.Error != nil {
		logger.Error(fmt.Sprintf("%v||%v", r1.Error.Error(), r2.Error.Error()))
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	data := make(map[string]interface{})
	data["count"] = articleCount
	data["recordList"] = archiveList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

// ===分页查看首页文章
type reqListArticles struct {
	Current int `form:"current" binding:"required"`
}

type tagDTO struct {
	TagId   int64  `json:"id"`
	TagName string `json:"tagName"`
}

type articleInfo struct {
	ID             int64     `json:"id"`
	ArticleCover   string    `json:"articleCover"`
	ArticleTitle   string    `json:"articleTitle"`
	ArticleContent string    `json:"articleContent"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
	IsTop          int       `json:"isTop"`
	Type           int       `json:"type"`
	CategoryId     int64     `json:"categoryId"`
	CategoryName   string    `json:"categoryName"`
	TagDTOList     []tagDTO  `json:"tagDTOList"`
}

func (a *Article) ListArticles(ctx *gin.Context) {
	var form reqListArticles
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验错误")
		return
	}
	db := common.GetGorm()
	offset := (form.Current - 1) * 10 //设定页面大小为10
	rows, err := db.Raw("select * from v_article_info limit ?,?", offset, 10).Rows()
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	articleInfoList := make([]articleInfo, 0)
	for rows.Next() {
		var t articleInfo
		_ = db.ScanRows(rows, &t)
		articleInfoList = append(articleInfoList, t)
	}
	// 获取文章标签
	for i, val := range articleInfoList {
		tagList := make([]tagDTO, 0)
		rows, err := db.Raw("select article_id,tt.id  as tag_id ,tag_name from t_article_tag ta join t_tag tt on ta.tag_id = tt.id where article_id=?", val.ID).Rows()
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error(err.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for rows.Next() {
			var t tagDTO
			_ = db.ScanRows(rows, &t)
			tagList = append(tagList, t)
		}
		articleInfoList[i].TagDTOList = tagList
	}
	Response(ctx, errorcode.Success, articleInfoList, true, "操作成功")
	return

}
func (a *Article) ListArticleBack(ctx *gin.Context)     {}
func (a *Article) SaveOrUpdateArticle(ctx *gin.Context) {}
func (a *Article) UpdateArticleTop(ctx *gin.Context)    {}
func (a *Article) UpdateArticleDelete(ctx *gin.Context) {}
func (a *Article) SaveArticleImages(ctx *gin.Context)   {}
func (a *Article) DeleteArticle(ctx *gin.Context)       {}
func (a *Article) GetArticleBackById(ctx *gin.Context)  {}

//======= 根据文章id获取文章详情也
type reqGetArticleById struct {
	ArticleId int64 `uri:"articleId" binding:"required"`
}

type lastNextArticle struct {
	ID           int64     `json:"id"`
	ArticleCover string    `json:"articleCover"`
	ArticleTitle string    `json:"articleTitle"`
	CreateTime   time.Time `json:"createTime"`
}

type respGetArticleById struct {
	ID                   int64             `json:"id"`
	ArticleCover         string            `json:"articleCover"`
	ArticleTitle         string            `json:"articleTitle"`
	ArticleContent       string            `json:"articleContent"`
	CreateTime           time.Time         `json:"createTime"`
	UpdateTime           time.Time         `json:"updateTime"`
	IsTop                int               `json:"isTop"`
	Type                 int               `json:"type"`
	CategoryId           int64             `json:"categoryId"`
	CategoryName         string            `json:"categoryName"`
	TagDTOList           []tagDTO          `json:"tagDTOList"`
	OriginalUrl          string            `json:"originalUrl"`
	ViewCount            int64             `json:"viewsCount"`
	LikeCount            int64             `json:"likeCount"`
	LastArticle          lastNextArticle   `json:"lastArticle"`
	NextArticle          lastNextArticle   `json:"nextArticle"`
	RecommendArticleList []lastNextArticle `json:"recommendArticleList"`
	NewestArticleList    []lastNextArticle `json:"newestArticleList"`
}

func (a *Article) GetArticleById(ctx *gin.Context) {
	var form reqGetArticleById
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验错误")
		return
	}
	db := common.GetGorm()
	data := respGetArticleById{}
	row, r1 := db.Raw("select * from v_article_info where id = ?", form.ArticleId).Rows()
	if r1 != nil && r1 != gorm.ErrRecordNotFound {
		logger.Error(r1.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r1 != nil {
		Response(ctx, errorcode.NotFoundResource, nil, false, "找不到对应文章")
		return
	}
	for row.Next() {
		_ = db.ScanRows(row, &data)
	}
	//查找文章对应的标签
	rows, r5 := db.Raw("select article_id,tt.id  as tag_id ,tag_name from t_article_tag ta join t_tag tt on ta.tag_id = tt.id where article_id=?", data.ID).Rows()
	if r5 != nil && r5 != gorm.ErrRecordNotFound {
		logger.Error(r5.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for rows.Next() {
		var t tagDTO
		_ = db.ScanRows(rows, &t)
		data.TagDTOList = append(data.TagDTOList, t)
	}
	// 找前一篇文章
	r2 := db.Model(&common.TArticle{}).Where("id < ? and is_delete=0", data.ID).Order("id DESC").First(&data.LastArticle)
	if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
		logger.Error(r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	//查找下一篇文章
	r3 := db.Model(&common.TArticle{}).Where("id > ? and is_delete=0", data.ID).Order("id ASC").First(&data.NextArticle)
	if r3.Error != nil && r3.Error != gorm.ErrRecordNotFound {
		logger.Error(r3.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 获取5篇最新的文章
	r4 := db.Model(&common.TArticle{}).Where("is_delete = 0").Order("create_time DESC").Limit(5).Find(&data.NewestArticleList)
	if r4.Error != nil && r4.Error != gorm.ErrRecordNotFound {
		logger.Error(r4.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 推荐3篇浏览数最高的文章
	r6 := db.Model(&common.TArticle{}).Where("is_delete = 0").Order("view_count DESC").Limit(3).Find(&data.RecommendArticleList)
	if r6.Error != nil && r6.Error != gorm.ErrRecordNotFound {
		logger.Error(r6.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, data, true, "操作成功")

}
func (a *Article) ListArticleByCondition(ctx *gin.Context) {}

type reqListArticleBySearch struct {
	Current  int    `form:"current" binding:"required"`
	KeyWords string `form:"keywords" binding:"required"`
}

type articleSearch struct {
	ID             int64  `json:"id"`
	ArticleTitle   string `json:"articleTitle"`
	ArticleContent string `json:"articleContent"`
	IsDelete       int    `json:"isDelete"`
	Status         int    `json:"status"`
}

func (a *Article) ListArticleBySearch(ctx *gin.Context) {
	var form reqListArticleBySearch
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数验证失败")
		return
	}
	db := common.GetGorm()
	var articleList []articleSearch
	r := db.Table("t_article").Where(fmt.Sprintf("article_title LIKE %q or article_content LIKE %q", "%"+form.KeyWords+"%", "%"+form.KeyWords+"%")).Find(&articleList)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		logger.Error(r.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	// 高亮搜索关键词
	html := "<span style='color:#f47466'>%s</span>"
	for i, val := range articleList {
		t := strings.Split(val.ArticleTitle, form.KeyWords)
		if len(t) >= 2 {
			articleList[i].ArticleTitle = ""
			count := 0
			for _, val := range t {
				articleList[i].ArticleTitle += val + fmt.Sprintf(html, form.KeyWords)
				count++
				if count+1 == len(t) {
					break
				}
			}
			articleList[i].ArticleTitle += t[len(t)-1]

		}

		t = strings.Split(val.ArticleContent, form.KeyWords)
		if len(t) >= 2 {
			articleList[i].ArticleContent = ""
			count := 0
			for _, val := range t {
				articleList[i].ArticleContent += val + fmt.Sprintf(html, form.KeyWords)
				count++
				if count+1 == len(t) {
					break
				}
			}
			articleList[i].ArticleContent += t[len(t)-1]
		}

	}
	Response(ctx, errorcode.Success, articleList, true, "操作成功")

}
func (a *Article) SaveArticleLike(ctx *gin.Context) {}

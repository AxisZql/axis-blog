package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"mime/multipart"
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

type reqListArticleBack struct {
	Current    int    `form:"current"`
	Size       int    `form:"size"`
	Keywords   string `form:"keywords"`
	CategoryId int64  `form:"categoryId"`
	Status     int    `form:"status"`
	TagId      int64  `form:"tagId"`
	Type       int    `form:"type"`
	IsDelete   int    `form:"isDelete"`
}
type listArticleBackInfo struct {
	ID           int64     `json:"id"`
	ArticleCover string    `json:"articleCover"`
	ArticleTitle string    `json:"articleTitle"`
	CategoryName string    `json:"categoryName"`
	CategoryId   int64     `json:"categoryId"`
	IsDelete     int       `json:"isDelete"`
	IsTop        int       `json:"isTop"`
	LikeCount    int64     `json:"likeCount"`
	Status       int       `json:"status"`
	CreateTime   time.Time `json:"createTime"`
	Type         int       `json:"type"`
	ViewsCount   int64     `json:"viewsCount"`
	TagDTOList   []struct {
		ID      int64  `json:"id"`
		TagName string `json:"tagName"`
	} `json:"tagDTOList"`
}

func (a *Article) ListArticleBack(ctx *gin.Context) {
	var form reqListArticleBack
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	db := common.GetGorm()
	var aTagList []common.TArticleTag
	var aIdList []int64
	te := false
	if form.TagId != 0 {
		te = true
		r1 := db.Where("tag_id = ?", form.TagId).Find(&aTagList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for _, val := range aTagList {
			aIdList = append(aIdList, val.ArticleId)
		}
	}
	condition := ""
	find := func(c int64, k string) {
		if c != 0 {
			condition += fmt.Sprintf("AND category_id = %d ", c)
		}
		if k != "" {
			condition += fmt.Sprintf("AND article_content LIKE %q ", "%"+k+"%")
		}
	}
	if form.Status != 0 || form.Type != 0 {
		if form.Status != 0 && form.Type != 0 {
			condition += fmt.Sprintf("status = %d AND type = %d ", form.Status, form.Type)
			find(form.CategoryId, form.Keywords)
		} else if form.Status != 0 {
			condition += fmt.Sprintf("status = %d ", form.Status)
			find(form.CategoryId, form.Keywords)
		} else {
			condition += fmt.Sprintf("type = %d ", form.Type)
			find(form.CategoryId, form.Keywords)
		}
	} else {
		find(form.CategoryId, form.Keywords)
		_condition := strings.Split(condition, " ")
		if len(_condition) >= 1 && _condition[0] == "AND" {
			_condition = _condition[1:]
		}
		condition = strings.Join(_condition, " ")
	}

	var aInfoList []listArticleBackInfo
	var count int64
	if te {
		if condition == "" {
			condition = "id IN ? "
		} else {
			condition += "AND id IN ? "
		}
		condition += fmt.Sprintf("AND is_delete = %d ", form.IsDelete)
		r2 := db.Model(&common.TArticle{}).Where(condition, aIdList).Count(&count)
		r2 = db.Model(&common.TArticle{}).Where(condition, aIdList).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&aInfoList)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		if condition == "" {
			condition += fmt.Sprintf("is_delete = %d ", form.IsDelete)
		} else {
			condition += fmt.Sprintf("AND is_delete = %d", form.IsDelete)
		}
		r2 := db.Model(&common.TArticle{}).Where(condition).Count(&count)
		r2 = db.Model(&common.TArticle{}).Where(condition).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&aInfoList)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}

	for i, val := range aInfoList {
		var at []common.TArticleTag
		r3 := db.Where("article_id = ?", val.ID).Find(&at)
		if r3.Error != nil && r3.Error != gorm.ErrRecordNotFound {
			logger.Error(r3.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for _, v := range at {
			var _tag struct {
				ID      int64  `json:"id"`
				TagName string `json:"tagName"`
			}
			r4 := db.Model(&common.TTag{}).Where("id = ?", v.TagId).Find(&_tag)
			if r4.Error != nil && r4.Error != gorm.ErrRecordNotFound {
				logger.Error(r4.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			aInfoList[i].TagDTOList = append(aInfoList[i].TagDTOList, _tag)
		}
		var category common.TCategory
		r4 := db.Model(&common.TCategory{}).Where("id = ?", val.CategoryId).Find(&category)
		if r4.Error != nil && r4.Error != gorm.ErrRecordNotFound {
			logger.Error(r4.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		aInfoList[i].CategoryName = category.CategoryName
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = aInfoList
	Response(ctx, errorcode.Success, data, true, "操作成功")

}

type reqSaveOrUpdateArticle struct {
	ID             int64    `json:"id" `
	ArticleContent string   `json:"articleContent" binding:"required"`
	ArticleCover   string   `json:"articleCover" `
	ArticleTitle   string   `json:"articleTitle" binding:"required"`
	CategoryName   string   `json:"categoryName" `
	IsTop          int      `json:"isTop" `
	OriginalUrl    string   `json:"originalUrl" `
	Status         int      `json:"status"`
	TagNameList    []string `json:"tagNameList"`
	Type           int      `json:"type" binding:"required"`
}

func (a *Article) SaveOrUpdateArticle(ctx *gin.Context) {
	var form reqSaveOrUpdateArticle
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	if form.Status == 0 {
		form.Status = 1
	}
	_session, _ := Store.Get(ctx.Request, "CurUser")
	aid := _session.Values["a_userid"]
	var ua common.TUserAuth
	var ui common.TUserInfo

	db := common.GetGorm()
	r0 := db.Where("id = ?", aid).First(&ua)
	if r0.Error != nil {
		logger.Error(r0.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	r0 = db.Where("id = ?", ua.UserInfoId).First(&ui)
	if r0.Error != nil {
		logger.Error(r0.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var category common.TCategory
	r1 := db.Where("category_name = ?", form.CategoryName).First(&category)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r1.Error == gorm.ErrRecordNotFound {
		category.CategoryName = form.CategoryName
		r1 = db.Model(&common.TCategory{}).Create(&category)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	var tagListId []int64
	for _, val := range form.TagNameList {
		var tag common.TTag
		r2 := db.Where("tag_name = ?", val).First(&tag)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		if r2.Error == gorm.ErrRecordNotFound {
			tag.TagName = val
			r2 = db.Model(&common.TTag{}).Create(&tag)
			if r2.Error != nil {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		tagListId = append(tagListId, tag.ID)
	}
	var article common.TArticle
	r := db.Where("id = ?", form.ID).First(&article)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r.Error == gorm.ErrRecordNotFound {
		article = common.TArticle{
			ArticleContent: form.ArticleContent,
			ArticleCover:   form.ArticleCover,
			ArticleTitle:   form.ArticleTitle,
			CategoryId:     category.ID,
			UserId:         ui.ID,
			IsTop:          form.IsTop,
			OriginalUrl:    form.OriginalUrl,
			Status:         form.Status,
			Type:           form.Type,
		}
		r3 := db.Model(&common.TArticle{}).Create(&article)
		if r3.Error != nil {
			logger.Error(r3.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	} else {
		article = common.TArticle{
			ArticleContent: form.ArticleContent,
			ArticleCover:   form.ArticleCover,
			ArticleTitle:   form.ArticleTitle,
			CategoryId:     category.ID,
			UserId:         ui.ID,
			IsTop:          form.IsTop,
			OriginalUrl:    form.OriginalUrl,
			Status:         form.Status,
			Type:           form.Type,
			UpdateTime:     time.Now(),
		}
		r3 := db.Model(&common.TArticle{}).Where("id = ?", form.ID).Updates(&article)
		if r3.Error != nil {
			logger.Error(r3.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}

	for _, val := range tagListId {
		var t common.TArticleTag
		r4 := db.Where("article_id = ? and tag_id = ?", article.ID, val).First(&t)
		if r4.Error == nil {
			continue
		} else if r4.Error == gorm.ErrRecordNotFound {
			t.ArticleId = article.ID
			t.TagId = val
			r4 = db.Model(&common.TArticleTag{}).Create(&t)
			if r4.Error != nil {
				logger.Error(r4.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		} else {
			logger.Error(r4.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	form.ID = article.ID
	Response(ctx, errorcode.Success, form, true, "操作成功")
}

type reqUpdateArticleTop struct {
	ID    int64 `json:"id" binding:"required"`
	IsTop int   `json:"isTop"`
}

func (a *Article) UpdateArticleTop(ctx *gin.Context) {
	var form reqUpdateArticleTop
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var article common.TArticle
	r1 := db.Where("id = ?", form.ID).First(&article)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	article.IsTop = form.IsTop
	article.UpdateTime = time.Now()
	r1 = db.Save(&article)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqUpdateArticleDelete struct {
	IdList   []int64 `json:"idList"`
	IsDelete int     `json:"isDelete"`
}

func (a *Article) UpdateArticleDelete(ctx *gin.Context) {
	var form reqUpdateArticleDelete
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	for _, val := range form.IdList {
		var a common.TArticle
		r1 := db.Where("id = ?", val).First(&a)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		a.IsDelete = form.IsDelete
		a.UpdateTime = time.Now()
		r1 = db.Save(&a)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqSaveArticleImages struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (a *Article) SaveArticleImages(ctx *gin.Context) {
	var form reqSaveArticleImages
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
	filePath := common.Conf.App.ArticleDir + fileName
	err := ctx.SaveUploadedFile(form.File, filePath)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	imgUrl := fmt.Sprintf("%s:%d/farticles/%s", common.Conf.App.HostName, common.Conf.App.Port, fileName)
	Response(ctx, errorcode.Fail, imgUrl, true, "操作成功")
}

func (a *Article) DeleteArticle(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	str := fmt.Sprintf("%v", string(data))
	var idList []int64
	err := json.Unmarshal([]byte(str), &idList)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	db := common.GetGorm()
	for _, val := range idList {
		r1 := db.Where("id = ?", val).Delete(&common.TArticle{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		r1 = db.Where("article_id = ?", val).Delete(&common.TArticleTag{})
		if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}

type reqGetArticleBackById struct {
	ArticleId int64 `uri:"articleId" binding:"required"`
}

func (a *Article) GetArticleBackById(ctx *gin.Context) {
	var form reqGetArticleBackById
	if err := ctx.ShouldBindUri(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	var articleInfo reqSaveOrUpdateArticle
	db := common.GetGorm()
	r1 := db.Table("v_article_info").Where("id = ?", form.ArticleId).Find(&articleInfo)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	var tags []common.TArticleTag
	r2 := db.Where("article_id = ?", form.ArticleId).Find(&tags)
	if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
		logger.Error(r2.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for _, val := range tags {
		var t common.TTag
		r3 := db.Where("id = ?", val.ArticleId).First(&t)
		if r3.Error != nil && r3.Error != gorm.ErrRecordNotFound {
			logger.Error(r3.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		articleInfo.TagNameList = append(articleInfo.TagNameList, t.TagName)
	}
	Response(ctx, errorcode.Success, articleInfo, true, "操作成功")
}

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

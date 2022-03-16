package cmd

import (
	"blog-server/common"
	ctrl "blog-server/service"
	util "blog-server/service/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api 服务", //短描述
	Long:  "api 服务", //长描述
	Run: func(cmd *cobra.Command, args []string) {
		AxAPi()
	},
}

func AxAPi() {
	// 运行环境初始化
	common.EnvInit()
	router := gin.Default()
	Routers(router)
	port := fmt.Sprintf(":%d", common.Conf.App.Port)
	server := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 32 << 20,
	}
	logger.Info(fmt.Sprintf("服务在%s端口启动成功", port))
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Routers 路由设置
func Routers(r *gin.Engine) {
	login := ctrl.Login{}
	userAuth := ctrl.UserAuth{}
	userInfo := ctrl.UserInfo{}
	tag := ctrl.Tag{}
	talk := ctrl.Talk{}
	role := ctrl.Role{}
	photo := ctrl.Photo{}
	resource := ctrl.Resource{}
	photoAlbum := ctrl.PhotoAlbum{}
	message := ctrl.Message{}
	page := ctrl.Page{}
	menus := ctrl.Menu{}
	loggerHandler := ctrl.Logger{}
	friendLink := ctrl.FriendLink{}
	comment := ctrl.Comment{}
	category := ctrl.Category{}
	blogInfo := ctrl.BlogInfo{}
	article := ctrl.Article{}
	admin := r.Group("/admin")
	users := r.Group("/users")
	home := r.Group("/home")

	r.GET("/", blogInfo.GetBlogHomeInfo)                                  //查看博客信息
	r.POST("/login", login.Login)                                         //用户登陆
	r.GET("/logout", login.LoginOut)                                      //用户注销
	r.POST("/register", userAuth.Register)                                //用户注册
	r.GET("/talks", talk.ListTalks)                                       //查看说说列表
	r.GET("/talk/:talkId", talk.GetTalkById)                              //根据id查看说说
	r.POST("/talks/:talkId/like", talk.SaveTalkLike)                      //点赞说说
	r.GET("/tags", tag.ListTags)                                          //查询标签列表
	r.GET("/photos/albums", photoAlbum.ListPhotoAlbum)                    //获取相册列表
	r.POST("/message", message.SaveMessage)                               //添加留言
	r.GET("/message", message.ListMessage)                                //查看留言列表
	r.GET("/links", friendLink.ListFriendLinks)                           //查看友链列表
	r.GET("/comments", comment.ListComment)                               //查询评论
	r.POST("/comments", comment.SaveComment)                              //添加评论
	r.GET("/comments/:commentId/replies", comment.ListRepliesByCommentId) //查询评论下的回复
	r.POST("/comments/:commentId/like", comment.SaveCommentLike)          //评论点赞
	r.GET("/categories", category.ListCategories)                         //查看分类列表
	r.GET("/about", blogInfo.GetAbout)                                    //查看关于我信息
	r.POST("/voice", blogInfo.SendVoice)                                  //上传语音信息
	r.POST("/report", blogInfo.Report)                                    //上传访客信息
	r.GET("/articles/archives", article.ListArchives)                     //文章归档列表
	r.GET("/articles", article.ListArticles)                              //查看首页文章
	r.GET("/articles/:articleId", article.GetArticleById)                 //更加id查看文章
	r.GET("/articles/condition", article.ListArticleByCondition)          //根据条件查询文章
	r.GET("/articles/search", article.ListArticleBySearch)                //搜索文章
	r.GET("/articles/:articleId/like", article.SaveArticleLike)           //点赞文章

	home.Use(util.Auth())
	{
		home.GET("/talks", talk.ListHomeTalks) //查看首页说说
	}

	r.GET("/admin", blogInfo.GetBlogBackInfo) //查询后台信息
	users.Use(util.Auth())
	{
		users.GET("/code", userAuth.SendEmailCode)       //发送邮箱验证码
		users.PUT("/password", userAuth.UpdatePassword)  //修改密码
		users.POST("/oauth/weibo", userAuth.WeiboLogin)  //微博登陆
		users.POST("/users/oauth/qq", userAuth.QQLogin)  //QQ登陆
		users.PUT("/info", userInfo.UpdateUserInfo)      //更新用户信息
		users.POST("/avatar", userInfo.UpdateUserAvatar) //更新用户头像
		users.POST("/email", userInfo.SaveUserEmail)     //绑定用户邮箱
	}

	admin.Use(util.Auth())
	{
		admin.GET("/users/area", userAuth.ListUserAreas)                         //获取用户区域分布
		admin.GET("/users", userAuth.ListUsers)                                  //查询用户后台列表
		admin.PUT("/users/password", userAuth.UpdateAdminPassword)               //修改管理员密码
		admin.PUT("/users/role")                                                 //修改用户角色
		admin.PUT("/admin/users/disable", userInfo.UpdateUserDisable)            //修改用户禁用状态
		admin.GET("/users/online", userInfo.ListOnlineUsers)                     //查看在线用户
		admin.DELETE("/user/:userInfoId/online", userInfo.RemoveOnlineUser)      //下线用户
		admin.POST("/talks/images", talk.SaveTalkImages)                         //上传说说图片
		admin.POST("/talks", talk.SaveOrUpdateTalk)                              //保存或者更新说说
		admin.DELETE("/talks", talk.DeleteTalks)                                 //删除说说
		admin.GET("/talks", talk.ListBackTalks)                                  //查看后台说说
		admin.GET("/talks/:talkId", talk.GetBackTalkById)                        //根据id查看后台说说
		admin.GET("/tags", tag.ListTagBack)                                      //查询后台标签列表
		admin.GET("/tags/search", tag.ListTagBySearch)                           //搜索文章标签
		admin.POST("/tags", tag.SaveOrUpdateTag)                                 //添加或者修改标签
		admin.DELETE("/tags", tag.DeleteTag)                                     //删除标签
		admin.GET("/users/role", role.ListUserRoles)                             //查询用户角色选项
		admin.GET("/roles", role.ListRoles)                                      //查询角色列表
		admin.POST("/role", role.SaveOrUpdateRole)                               //保存或更新角色
		admin.DELETE("/roles", role.DeleteRoles)                                 //删除角色
		admin.GET("/resources", resource.ListResources)                          //查看资源列表
		admin.DELETE("/resources/:resourceId", resource.DeleteResource)          //删除资源
		admin.POST("/resources", resource.SaveOrUpdateResource)                  //新增或者修改资源
		admin.GET("/role/resources", resource.ListResourceOption)                //查看角色资源选项
		admin.GET("/photos", photo.ListPhotos)                                   //根据相册id获取照片列表
		admin.PUT("/photos", photo.UpdatePhoto)                                  //更新照片信息
		admin.POST("/photos", photo.SavePhoto)                                   //保存照片
		admin.PUT("/photos/album", photo.UpdatePhotoAlbum)                       //移动照片相册
		admin.PUT("/photos/delete", photo.UpdatePhotoDelete)                     //更新照片删除状态
		admin.DELETE("/photos", photo.DeletePhotos)                              //删除照片
		admin.GET("/:albumId/photos", photo.ListPhotoByAlbumId)                  //根据相册id查看照片类别
		admin.POST("/photos/albums/cover", photoAlbum.SavePhotoAlbumCover)       //上传相册封面
		admin.POST("/photos/albums", photoAlbum.SaveOrUpdatePhotoAlbum)          //保存或者更新相册
		admin.GET("/photos/albums", photoAlbum.ListPhotoAlbumBack)               //查看后台相册列表
		admin.GET("/photos/albums/info", photoAlbum.ListPhotoAlbumBackInfo)      //获取后台相册相关信息
		admin.GET("/photos/:albumsId/info", photoAlbum.GetPhotoAlbumBackById)    //根据id获取后台相册信息
		admin.DELETE("/photos/albums/:albumId", photoAlbum.DeletePhotoAlbumById) //根据相册id删除相册
		admin.DELETE("/pages/:pageId", page.DeletePage)                          //根据页面id删除页面
		admin.POST("/pages", page.SaveOrUpdatePage)                              //保存或者更新页面
		admin.GET("/pages", page.ListPages)                                      //获取页面列表
		admin.GET("/message", message.ListMessageBack)                           //查看后台留言列表
		admin.POST("/message/review", message.UpdateMessageReview)               //审核留言
		admin.DELETE("/messages", message.DeleteMessage)                         //删除留言
		admin.GET("/menus", menus.ListMenus)                                     //查看菜单列表
		admin.POST("/menus", menus.SaveOrUpdateMenu)                             //新增或者修改菜单
		admin.DELETE("/menus/:menuId", menus.DeleteMenu)                         //删除菜单
		admin.GET("/role/menus", menus.ListUserMenus)                            //查看当前用户菜单
		admin.GET("/operation/logs", loggerHandler.ListOperationLogs)            //查看操作日志
		admin.DELETE("/operation/logs", loggerHandler.DeleteOperationLogs)       //删除操作日志
		admin.GET("/links", friendLink.ListFriendLinksBack)                      //查看后台友链列表
		admin.POST("/links", friendLink.SaveOrUpdateFriendLink)                  //保存或修改友链
		admin.DELETE("/links", friendLink.DeleteFriendLink)                      //删除友链
		admin.PUT("/comments/reviews", comment.UpdateCommentReview)              //评论审核
		admin.DELETE("/comments", comment.DeleteComment)                         //删除评论
		admin.GET("/comments", comment.ListCommentBack)                          //查询后台评论
		admin.GET("/categories", category.ListCategoriesBack)                    //查看后台分类列表
		admin.GET("/categories/search", category.ListCategoriesBySearch)         //搜索文章分类
		admin.POST("/categories", category.SaveOrUpdateCategory)                 //添加或者修改分类
		admin.DELETE("/categories", category.DeleteCategories)                   //删除分类
		admin.POST("/config/images", blogInfo.SavePhotoAlbumCover)               //上传博客配置图片
		admin.PUT("/website/config", blogInfo.UpdateWebsiteConfig)               //更新网站配置
		admin.PUT("/about", blogInfo.UpdateAbout)                                //修改关于我信息
		admin.GET("/articles", article.ListArticleBack)                          //查看后台文章
		admin.POST("/articles", article.SaveOrUpdateArticle)                     //添加或者修改文章
		admin.PUT("/articles/top", article.UpdateArticleTop)                     //修改文章置顶
		admin.PUT("/articles", article.UpdateArticleDelete)                      //恢复或删除文章
		admin.POST("/articles/images", article.SaveArticleImages)                //上传文章图片
		admin.DELETE("/articles", article.DeleteArticle)                         //物理删除文章
		admin.GET("/articles/:articleId", article.GetArticleBackById)

	}

}

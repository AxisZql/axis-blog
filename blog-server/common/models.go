package common

import (
	"time"
)

/**
* @author:AxisZql
* @date:2022-3-15 8:50 AM
* @desc:数据库模型定义
 */

// TArticle 文章表
type TArticle struct {
	ID             int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId         int64     `json:"user_id,omitempty" remark:"作者id" gorm:"type:bigint;not null"`
	CategoryId     int64     `json:"category_id,omitempty" remark:"文章分类id" gorm:"type:bigint"`
	ArticleCover   string    `json:"article_cover,omitempty" remark:"文章封面" gorm:"type:varchar(1024)"`
	ArticleTitle   string    `json:"article_title,omitempty" remark:"标题" gorm:"type:varchar(50);not null"`
	ArticleContent string    `json:"article_content,omitempty" remark:"内容" gorm:"type:longtext;not null"`
	Type           int       `json:"type,omitempty" remark:"文章类型 1原创 2转载 3翻译" gorm:"type:tinyint;default:0;not null"`
	OriginalUrl    string    `json:"original_url,omitempty" remark:"原文链接" gorm:"type:varchar(255)"`
	ViewCount      int64     `json:"view_count,omitempty" remark:"文章访问量" gorm:"type:bigint;default:0;not null"`
	LikeCount      int64     `json:"like_count,omitempty" remark:"文章点赞数" gorm:"type:bigint;default:0"`
	IsTop          int       `json:"is_top,omitempty" remark:"是否置顶 0否 1是" gorm:"type:tinyint;default:0;not null"`
	IsDelete       int       `json:"is_delete,omitempty" remark:"是否删除  0否 1是" gorm:"type:tinyint;default:0;not null"`
	Status         int       `json:"status,omitempty" remark:"状态值 1公开 2私密 3评论可见" gorm:"type:tinyint;default:1;not null"`
	CreateTime     time.Time `json:"create_time,omitempty" remark:"发表时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime     time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TLike 点赞表 {"articleLikeSet":[],"commentLikeSet:[],"talkLikeSet":[]}
type TLike struct {
	ID       int64  `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId   int64  `json:"user_id,omitempty" remark:"用户id" gorm:"type:bigint;not null"`
	LikeItem string `json:"article_like,omitempty" remark:"用户点赞文章id，评论id，说说id数组" gorm:"type:text"`
}

// TCategory 分类表
type TCategory struct {
	ID           int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	CategoryName string    `json:"category_name,omitempty" remark:"分类名" gorm:"type:varchar(20);not null"`
	CreateTime   time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime   time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TTag 标签表
type TTag struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	TagName    string    `json:"tag_name,omitempty" remark:"标签名" gorm:"type:varchar(20);not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TArticleTag 文章所属标签表
type TArticleTag struct {
	ID        int64 `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	ArticleId int64 `json:"article_id" remark:"文章id" gorm:"type:bigint;not null;index"`
	TagId     int64 `json:"tag_id" remark:"标签id" gorm:"type:bigint;not null;index"`
}

// TChatRecord 聊天记录表
type TChatRecord struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId     int64     `json:"user_id,omitempty" remark:"用户id(匿名用户没有id)" gorm:"type:bigint"`
	Nickname   string    `json:"nickname,omitempty" remark:"用户昵称" gorm:"type:varchar(50);not null"`
	Avatar     string    `json:"avatar,omitempty" remark:"用户头像" gorm:"type:varchar(255);not null"`
	Content    string    `json:"content,omitempty" remark:"聊天内容" gorm:"type:varchar(1024);not null"`
	IpAddress  string    `json:"ip_address,omitempty" remark:"ip地址" gorm:"type:varchar(50);not null"`
	IpSource   string    `json:"ip_source,omitempty" remark:"ip来源" gorm:"type:varchar(255);not null"`
	Type       int       `json:"type,omitempty" remark:"类型" gorm:"type:tinyint;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TComment 评论表
type TComment struct {
	ID             int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId         int64     `json:"user_id" remark:"用户id" gorm:"type:bigint;not null;index"`
	TopicId        int64     `json:"topic_id,omitempty" remark:"评论主题id" gorm:"type:bigint"`
	CommentContent string    `json:"comment_content" remark:"评论内容" gorm:"type:text;not null"`
	ReplyUserId    int64     `json:"reply_user_id" remark:"回复用户id" gorm:"type:bigint"`
	ParentId       int64     `json:"parent_id,omitempty" remark:"父评论id" gorm:"type:bigint;index"`
	Type           int       `json:"type,omitempty" remark:"评论类型 1.文章 2.友链 3.说说" gorm:"type:tinyint"`
	IsDelete       int       `json:"is_delete,omitempty" remark:"是否删除 0否 1是" gorm:"type:tinyint;default:0"`
	IsReview       int       `json:"is_review,omitempty" remark:"是否审核 0否 1是" gorm:"type:tinyint;default:1"`
	CreateTime     time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime     time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TFriendLink 友链表
type TFriendLink struct {
	ID          int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	LinkName    string    `json:"link_name,omitempty" remark:"链接名" gorm:"type:varchar(20);not null"`
	LinkAvatar  string    `json:"link_avatar,omitempty" remark:"链接头像" gorm:"type:varchar(255);not null"`
	LinkAddress string    `json:"link_address,omitempty" remark:"链接地址" gorm:"type:varchar(50);not null"`
	LinkIntro   string    `json:"link_intro,omitempty" remark:"链接介绍" gorm:"type:varchar(100);not null"`
	CreateTime  time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime  time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TMenu 菜单表
type TMenu struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	Name       string    `json:"name,omitempty" remark:"菜单名" gorm:"type:varchar(20);not null"`
	Path       string    `json:"path,omitempty" remark:"菜单路径" gorm:"type:varchar(50);not null"`
	Component  string    `json:"component,omitempty" remark:"组件" gorm:"type:varchar(50);not null"`
	Icon       string    `json:"icon,omitempty" remark:"菜单icon" gorm:"type:varchar(50);not null"`
	OrderNum   int64     `json:"order_num,omitempty" remark:"排序" gorm:"type:bigint;not null"`
	ParentId   int64     `json:"parent_id,omitempty" remark:"父菜单id" gorm:"type:bigint"`
	IsHidden   int       `json:"is_hidden,omitempty" remark:"是否隐藏 0否 1是" gorm:"type:tinyint;default:0;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TMessage 留言表
type TMessage struct {
	ID             int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	Nickname       string    `json:"nickname,omitempty" remark:"用户昵称" gorm:"type:varchar(50);not null"`
	Avatar         string    `json:"avatar,omitempty" remark:"用户头像" gorm:"type:varchar(255);not null"`
	MessageContent string    `json:"message_content,omitempty" remark:"留言内容" gorm:"type:varchar(255);not null"`
	IpAddress      string    `json:"ip_address,omitempty" remark:"ip地址" gorm:"type:varchar(50);not null"`
	IpSource       string    `json:"ip_source,omitempty" remark:"ip来源" gorm:"type:varchar(255);not null"`
	Speed          int       `json:"speed,omitempty" remark:"弹幕速度" gorm:"type:tinyint;not null"`
	CreateTime     time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime     time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TOperationLog 操作日志表
type TOperationLog struct {
	ID            int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	OptModule     string    `json:"opt_module,omitempty" remark:"操作模块" gorm:"type:varchar(20);not null"`
	OptType       string    `json:"opt_type,omitempty" remark:"操作类型" gorm:"type:varchar(20);not null"`
	OptUrl        string    `json:"opt_url,omitempty" remark:"操作url" gorm:"type:varchar(255);not null"`
	OptMethod     string    `json:"opt_method,omitempty" remark:"操作方法" gorm:"type:varchar(255);not null"`
	OptDesc       string    `json:"opt_desc,omitempty" remark:"操作描述" gorm:"type:varchar(255);not null"`
	RequestParam  string    `json:"request_param,omitempty" remark:"请求参数" gorm:"type:longtext;not null"`
	RequestMethod string    `json:"request_method,omitempty" remark:"请求方法" gorm:"type:varchar(20);not null"`
	ResponseData  string    `json:"response_data,omitempty" remark:"响应数据" gorm:"type:longtext;not null"`
	UserId        int64     `json:"user_id,omitempty" remark:"用户id" gorm:"type:bigint;not null"`
	Nickname      string    `json:"nickname,omitempty" remark:"用户昵称" gorm:"type:varchar(50);not null"`
	IpAddress     string    `json:"ip_address,omitempty" remark:"ip地址" gorm:"type:varchar(50);not null"`
	IpSource      string    `json:"ip_source,omitempty" remark:"ip来源" gorm:"type:varchar(255);not null"`
	CreateTime    time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime    time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TPage 页面表
type TPage struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	PageName   string    `json:"page_name,omitempty" remark:"页面名称" gorm:"type:varchar(10);not null"`
	PageLabel  string    `json:"page_label,omitempty" remark:"页面标签" gorm:"type:varchar(20)"`
	PageCover  string    `json:"page_cover,omitempty" remark:"页面封面" gorm:"type:varchar(255);not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TPhoto 照片表
type TPhoto struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	AlbumId    int64     `json:"album_id" remark:"相册id" gorm:"type:bigint;not null"`
	PhotoName  string    `json:"photo_name,omitempty" remark:"照片名" gorm:"type:varchar(20);not null"`
	PhotoDesc  string    `json:"photo_desc,omitempty" remark:"照片描述" gorm:"type:varchar(50)"`
	PhotoSrc   string    `json:"photo_src" remark:"照片地址" gorm:"type:varchar(255);not null"`
	IsDelete   int       `json:"is_delete,omitempty" remark:"是否删除 0否 1是" gorm:"type:tinyint;default:0;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TPhotoAlbum 相册表
type TPhotoAlbum struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	AlbumName  string    `json:"album_name,omitempty" remark:"相册名" gorm:"type:varchar(20);not null"`
	AlbumDesc  string    `json:"album_desc,omitempty" remark:"相册描述" gorm:"type:varchar(50);not null"`
	AlbumCover string    `json:"album_cover,omitempty" remark:"相册封面" gorm:"type:varchar(255);not null"`
	IsDelete   int       `json:"is_delete,omitempty" remark:"是否删除" gorm:"type:tinyint;default:0;not null"`
	Status     int       `json:"status,omitempty" remark:"状态值 1公开 2私密" gorm:"type:tinyint;default:1;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TResource 资源权限表
type TResource struct {
	ID            int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	ResourceName  string    `json:"resource_name,omitempty" remark:"资源名" gorm:"type:varchar(50);not null"`
	Url           string    `json:"url,omitempty" remark:"权限路径" gorm:"type:varchar(255)"`
	RequestMethod string    `json:"request_method" remark:"请求方式" gorm:"type:varchar(10)"`
	ParentId      int64     `json:"parent_id,omitempty" remark:"父权限id" gorm:"type:bigint"`
	IsAnonymous   int       `json:"is_anonymous,omitempty" remark:"是否匿名 0否 1是" gorm:"type:tinyint;default:0;not null"`
	CreateTime    time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime    time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TRole 角色表
type TRole struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	RoleName   string    `json:"role_name,omitempty" remark:"角色名" gorm:"type:varchar(20);not null"`
	RoleLabel  string    `json:"role_label,omitempty" remark:"角色描述" gorm:"type:varchar(50);not null"`
	IsDisable  int       `json:"is_disable,omitempty" remark:"是否禁用 0否 1是" gorm:"type:tinyint;default:0;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TRoleMenu 角色使用的菜单表
type TRoleMenu struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	RoleId     int64     `json:"role_id,omitempty" remark:"角色id" gorm:"type:bigint"`
	MenuId     int64     `json:"menu_id,omitempty" remark:"菜单id" gorm:"type:bigint"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TRoleResource 角色资源表
type TRoleResource struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	RoleId     int64     `json:"role_id,omitempty" remark:"角色id" gorm:"type:bigint"`
	ResourceId int64     `json:"resource_id,omitempty" remark:"权限id" gorm:"type:bigint"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TTalk 说说表
type TTalk struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId     int64     `json:"user_id,omitempty" remark:"用户id" gorm:"type:bigint;not null"`
	Content    string    `json:"content,omitempty" remark:"说说内容" gorm:"type:varchar(2000);not null"`
	Images     string    `json:"images,omitempty" remark:"说说图片" gorm:"type:varchar(2500)"`
	IsTop      int       `json:"is_top,omitempty" remark:"是否置顶 0否 1是" gorm:"type:tinyint;default 0"`
	Status     int       `json:"status,omitempty" remark:"状态 1公开 2私密" gorm:"type:tinyint;default:1;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TUniqueView 网站访问量表
type TUniqueView struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	ViewsCount int64     `json:"views_count,omitempty" remark:"访问量" gorm:"type:bigint;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TUserAuth 用户登陆身份验证记录表
type TUserAuth struct {
	ID            int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserInfoId    int64     `json:"user_info_id,omitempty" remark:"用户信息id" gorm:"type:bigint;not null"`
	Username      string    `json:"username,omitempty" remark:"用户名" gorm:"type:varchar(50);unique_index;not null"`
	Password      string    `json:"password,omitempty" remark:"密码" gorm:"type:varchar(255);not null"`
	LoginType     int       `json:"login_type,omitempty" remark:"登陆类型 1账号密码 2QQ 3微博" gorm:"type:tinyint;not null"`
	LastLoginTime time.Time `json:"last_login_time,omitempty" remark:"上次登陆时间" gorm:"type:datetime"`
	UserAgent     string    `json:"user_agent,omitempty" remark:"浏览器请求头UserAgent包含操作系统和浏览器信息" gorm:"type:varchar(255);not null"`
	IpAddress     string    `json:"ip_address,omitempty" remark:"ip地址" gorm:"type:varchar(50);not null"`
	IpSource      string    `json:"ip_source,omitempty" remark:"ip来源" gorm:"type:varchar(255);not null"`
	CreateTime    time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime    time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TUserInfo 用户信息表
type TUserInfo struct {
	ID         int64     `json:"id,omitempty" remark:"自增id(这是userid)" gorm:"primary_key"`
	Email      string    `json:"email,omitempty" remark:"邮箱号" gorm:"type:varchar(50);unique_index"`
	Nickname   string    `json:"nickname,omitempty" remark:"用户昵称" gorm:"type:varchar(50);not null"`
	Avatar     string    `json:"avatar,omitempty" remark:"用户头像" gorm:"type:varchar(255);not null"`
	Intro      string    `json:"intro,omitempty" remark:"用户简介" gorm:"type:varchar(255)"`
	WebSite    string    `json:"web_site" remark:"个人网站" gorm:"type:varchar(255)"`
	IsDisable  int       `json:"is_disable,omitempty" remark:"是否禁用 0否 1是" gorm:"type:tinyint;default:0;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TUserRole 用户权限表
type TUserRole struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	UserId     int64     `json:"user_id,omitempty" remark:"用户id" gorm:"type:bigint;not null"`
	RoleId     int64     `json:"role_id,omitempty" remark:"角色id" gorm:"type:bigint;not null"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

// TWebsiteConfig 网站配置表
type TWebsiteConfig struct {
	ID         int64     `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	Config     string    `json:"config,omitempty" remark:"配置信息" gorm:"type:varchar(1000)"`
	CreateTime time.Time `json:"create_time,omitempty" remark:"创建时间" gorm:"type:datetime;default:current_timestamp;not null"`
	UpdateTime time.Time `json:"update_time,omitempty" remark:"更新时间" gorm:"type:datetime;default:current_timestamp;not null"`
}

//==== 视图区

type VUserMenu struct {
	UserId    int64  `json:"user_id"`
	MenuId    int64  `json:"menu_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	OrderNum  string `json:"order_num"`
	ParentId  int64  `json:"parent_id"`
	IsHidden  int    `json:"is_hidden"`
}

type VCateGoryCount struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
	Count        int64  `json:"count"`
}

type VArticleStatistics struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

type VComment struct {
	ID             int64     `json:"id"`
	ParentId       int64     `json:"parent_id"`
	UserId         int64     `json:"user_id"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	WebSite        string    `json:"web_site"`
	ReplyUserId    int64     `json:"reply_user_id"`
	CommentContent string    `json:"comment_content"`
	TopicId        int64     `json:"topic_id"`
	Type           int       `json:"type"`
	ReplyNickname  string    `json:"reply_nickname"`
	ReplyWebSite   string    `json:"reply_web_site"`
	CreateTime     time.Time `json:"create_time"`
}

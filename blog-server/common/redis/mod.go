package redis

import "time"

// UserCodeKey 验证码
const UserCodeKey string = "code:%v"

// CodeExpireTime 验证码过期时间
const CodeExpireTime = 5 * time.Minute

// BlogViewsCount 博客浏览量
const BlogViewsCount string = "blog_views_count"

// ArticleViewsCount 文章浏览量
const ArticleViewsCount string = "article_views_count"

// ArticleLikeCount 文章点赞量
const ArticleLikeCount string = "article_like_count"

// TalkLikeCount 说说点赞量
const TalkLikeCount string = "talk_like_count"

// CommentLikeCount 评论点赞量
const CommentLikeCount string = "comment_like_count"

// UserLike 用户点赞数据
const UserLike string = "user_like:%v"
const ExpireUserLike = 5 * time.Second

// WebsiteConfig 网站配置
const WebsiteConfig string = "website_config"

// UserArea 用户地区
const UserArea string = "user_area"

// VisitorArea 访客地区
const VisitorArea string = "visitor_area"

// PageCover 页面封面
const PageCover string = "page_cover"

// ABOUT 关于我信息
const ABOUT string = "about"

// UniqueVisitor 访客数量
const UniqueVisitor string = "unique_visitor"

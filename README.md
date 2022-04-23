# axis-blog 介绍

<p align=center>
  <a href="https://magic-nubmer-tool.top"></a>
  <img src="https://static.magic-number-tool.top/config/d700e8ea9d33fc837b568fdfb454c8c5.jpeg" alt="axiszql blog" style="border-radius:50%;width=150px;height:150px">
</p>

<p align=center>
  基于Gin + Vue开发的前后端分离的博客，前端由<a href="https://github.com/X1192176811">风宇</a>巨佬开发，本人用Go语言对该博客进行后端全面改写。
</p>

<p align="center">
  <a target="_blank" href="https://github.com/AxisZql/axis-blog">
    <img src="https://img.shields.io/apm/l/vim-mode"/>
    <img src="https://img.shields.io/badge/redis-6.2.6-red&logo=Redis"/>
    <img src="https://img.shields.io/badge/GO-1.17-00acd7?logo=Go&logoColor=00acd7"/>
    <img src="https://img.shields.io/badge/Gin-1.7.7-blue"/>
    <img src="https://img.shields.io/badge/MySQL-8.0+-blue"/>
    <img src="https://img.shields.io/badge/WebSocket-laster-brightgreen"/>
     <img src="https://img.shields.io/badge/vue-2.5.17-green"/>
  </a>
</p>

[项目预览地址](#预览地址) |[目录结构](#目录结构)｜[项目特点](#项目特点)｜[技术介绍](#技术介绍)|[运行环境](#运行环境)｜[开发环境](#开发环境)|[项目截图](#项目截图)｜[部署方法](#部署方法)



## 预览地址

**项目访问地址:**  [https://magic-number-tool.top]([https://magic-number-tool.top])

**Github地址:** [https://github.com/AxisZql/axis-blog](https://github.com/AxisZql/axis-blog)

## 目录结构

```shell
.
├── Dockerfile
├── cmd
│   ├── axapi.go # 接口定义
│   ├── inti.go # cmd包下全局变量
│   └── root.go # 服务启动
├── common
│   ├── conf.go # 读取config.toml中的配置
│   ├── curd.go  # 数据库操作抽象
│   ├── curd_interface.go # 数据库操作接口定义
│   ├── database.go # 视图、表等初始化相关
│   ├── errorcode # 业务错误码定义
│   ├── init.go # common包下全局变量
│   ├── log.go # 日志记录器初始化
│   ├── models.go # 数据库模型定义
│   ├── redis.go # redis相关
│   ├── rediskey # redis 相关key定义
│   ├── tools # 工具包
│   │   ├── senitive_word.go # 敏感词检测
│   │   └── senitive_word_test.go # 敏感词检测测试
│   ├── utils.go # 工具类
│   └── utils_test.go # 每个工具方法的测试
├── config.toml # 配置文件
├── controllers # 所有RESTful API处理方法的接口定义
│   ├── article.go
│   ├── blog_info.go
│   ├── category.go
│   ├── comment.go
│   ├── friend_link.go
│   ├── logger.go
│   ├── login.go
│   ├── menu.go
│   ├── message.go
│   ├── page.go
│   ├── photo.go
│   ├── photo_album.go
│   ├── resource.go
│   ├── role.go
│   ├── tag.go
│   ├── talk.go
│   ├── user_auth.go
│   ├── user_info.go
│   └── websocket.go
├── deployment # 部署相关
│   ├── blog-mysql8.sql #数据库备份
│   └── nginx # nginx配置相关
│       ├── Dockerfile # nginx 容器的Dockerfile文件
│       ├── cert # https 证书存放位置
│       ├── frontend # 前端打包文件存放位置
│       └── nginx.conf # nginx配置文件
├── docker-compose.yml # 项目部署的docker-compose文件
├── go.mod
├── go.sum
├── main.go # 项目入口
├── sensitive-words.txt # 敏感词库
├── service # 所有RESTful API处理方法的接口实现
    ├── article.go
    ├── blog_info.go
    ├── category.go
    ├── comment.go
    ├── friend_link.go
    ├── init.go
    ├── logger.go
    ├── login.go
    ├── menu.go
    ├── message.go
    ├── page.go
    ├── photo.go
    ├── photo_album.go
    ├── resource.go
    ├── role.go
    ├── tag.go
    ├── talk.go
    ├── user_auth.go
    ├── user_info.go
    ├── utils.go
    └── websocket.go

```

## 项目特点

- 文章发布采用Markdown编辑器，高效简单
- 支持20余种表情进行评论和回复
- 使用WebSocket实现支持多人在线实时聊天
- 支持准确率极好的敏感词检测，确保聊天室、评论区、留言区出现不当言论
- 支持文章搜索关键词高亮操作
- 支持在线播放网易云热门音乐
- 留言采用弹幕形式
- 项目前后端分离部署
- 支持文章代码高亮功能

## 技术介绍

**前端：** vue + vuex + vue-router + axios + vuetify + element + echarts

**后端：** Gin + nginx + docker + gorm +Mysql + Redis + Websocket

## 运行环境

**服务器：** 阿里云2核2G CentOS8.2（最低配置）

## 开发环境

|开发工具|说明|
|-|-|
|GoLand|Go开发工具IDE|
|VSCode|Vue开发工具IDE|
|GoLand|GoLand自带MySQL远程连接工具|
|Another Redis Desktop Manager|Redis远程连接工具|
|Tabby|Linux远程连接工具|
|Transmit|Linux文件上传工具|
|Mac OS|操作系统|

|开发环境|版本|
|-|-|
|Go|1.17.2|
|MySQL|8.0+|
|Redis|6.2.6|

## 项目截图

![image.png](https://static.magic-number-tool.top/articles/a3bd02862df8d7f5adcf2d4dab63d36f.png)

![image.png](https://static.magic-number-tool.top/articles/3d816006bb2756549246e9844d284c22.png)



![image.png](https://static.magic-number-tool.top/articles/9a6652f4d549550109184be38ff086a5.png)

![image.png](https://static.magic-number-tool.top/articles/4ae448e10fd42f2654833595e5e3404d.png)

## 部署方法

> 先创建一个数据库然后把blog-mysql8.sql中的所有数据导入该新建数据库中

**Nginx 配置：**

```nginx

server {
           listen  443 ssl;
           server_name  你的前台域名;
           ssl on;
           ssl_certificate    你的前台域名pem证书位置;
           ssl_certificate_key  你的前台域名key证书位置;
           ssl_session_timeout 5m;
           ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
           ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
           ssl_prefer_server_ciphers on;

       location / {
            root   /usr/share/nginx/html/vue/blog;
            index  index.html index.htm;
            try_files $uri $uri/ /index.html;
        }

	   location ^~ /api/ {
            proxy_pass http://axisblog:9080/;
	        proxy_set_header   Host             $host;
            proxy_set_header   X-Real-IP        $remote_addr;
            proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
        }
    }


server {
        listen       80;
        server_name  你的前台域名;
        rewrite ^(.*)$	https://$host$1	permanent;
    }


server {
        listen  443 ssl;
        server_name  你的后台子域名;
        ssl on;
        ssl_certificate    你的后台子域名pem证书位置;
        ssl_certificate_key 你的后台子域名key证书位置;
        ssl_session_timeout 5m;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_prefer_server_ciphers on;

       location / {
                   root   /usr/share/nginx/html/vue/admin;
                   index  index.html index.htm;
                   try_files $uri $uri/ /index.html;
               }

       	location ^~ /api/ {
                   proxy_pass http://axisblog:9080/;
       	           proxy_set_header   Host             $host;
                   proxy_set_header   X-Real-IP        $remote_addr;
                   proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
               }
    }


server {
        listen       80;
        server_name  你的后台子域名;
        rewrite ^(.*)$ https://$host$1 permanent;
    }


server {
        listen  443 ssl;
        server_name  你的静态文件服务子域名;
        ssl on;
        ssl_certificate    你的静态文件服务子域名pem证书位置;;
        ssl_certificate_key   你的静态文件服务子域名key证书位置;;
        ssl_session_timeout 5m;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_prefer_server_ciphers on;
        #charset koi8-r;
        location / {
              root /usr/share/upload/static/;
        }
}

server {
        listen       80;
        server_name  你的静态文件服务子域名;
        rewrite ^(.*)$ https://$host$1 permanent;

    }
#


server {

        listen  443 ssl;
        server_name  你的websocket服务子域名;

        ssl on;
        ssl_certificate    你的websocket服务子域名pem证书位置;
        ssl_certificate_key 你的websocket服务子域名key证书位置;
        ssl_session_timeout 5m;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_prefer_server_ciphers on;

        location / {
          proxy_pass http://axisblog:9080/ws;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "Upgrade";
          proxy_set_header Host $host:$server_port;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto $scheme;
       }

    }

server {
        listen       80;
        server_name  你的websocket服务子域名;
        rewrite ^(.*)$	https://$host$1	permanent;

    }



```

**服务启动**:

1. 构建镜像

   ```shell
   docker-compose build
   ```

2. 启动容器

   ```shell
   docker-compose up -d
   ```

   

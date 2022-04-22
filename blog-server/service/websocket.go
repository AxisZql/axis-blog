package service

import (
	"blog-server/common"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type MyWebSocket struct{}

type Client struct {
	ID         string
	IpAddress  string
	IpSource   string
	UserId     interface{}
	Socket     *websocket.Conn
	Send       chan []byte
	Start      time.Time
	ExpireTime time.Duration // 一段时间没有接收到心跳则过期
}

type ClientManager struct {
	Clients    map[string]*Client // 记录在线用户
	Broadcast  chan []byte
	Register   chan *Client // 触发新用户登陆
	UnRegister chan *Client // 触发用户退出
}

type WsMessage struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

var Manager *ClientManager

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-Manager.Register:
			Manager.Clients[conn.ID] = conn
			// 如果有新用户连接则发送最近聊天记录和在线人数给他
			count := len(Manager.Clients)
			Manager.InitSend(conn, count)
		}
	}
}

// Quit 离线用户触发删除
func (manager *ClientManager) Quit() {
	for {
		select {
		case conn := <-Manager.UnRegister:
			delete(Manager.Clients, conn.ID)
			// 给客户端刷新在线人数
			resp, _ := json.Marshal(&WsMessage{Type: 1, Data: len(Manager.Clients)})
			manager.Broadcast <- resp
		}
	}
}

func (manager *ClientManager) InitSend(cur *Client, count int) {
	resp, _ := json.Marshal(&WsMessage{Type: 1, Data: count})
	Manager.Broadcast <- resp
	db := common.GetGorm()
	// 获取消息历史记录(12条最新)
	var chatList []common.TChatRecord
	r1 := db.Limit(12).Order("create_time DESC").Find(&chatList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		cur.Send <- []byte("系统异常")
	}
	// 按照时间正序排序
	sort.Slice(chatList, func(i, j int) bool {
		return chatList[i].CreateTime.Unix() < chatList[j].CreateTime.Unix()
	})
	_data := make(map[string]interface{})
	_data["chatRecordList"] = chatList
	_data["ipAddress"] = cur.IpAddress
	_data["ipSource"] = cur.IpSource
	resp, _ = json.Marshal(&WsMessage{Type: 2, Data: _data})
	cur.Send <- resp
}

// BroadcastSend 群发消息
func (manager *ClientManager) BroadcastSend() {
	for {
		select {
		// 只要有一方发消息就广播
		case msg := <-Manager.Broadcast:
			for _, conn := range Manager.Clients {
				conn.Send <- msg
			}
		}
	}
}

// Check 实时监测过期
func (c *Client) Check() {
	for {
		now := time.Now()
		var duration = now.Sub(c.Start)
		if duration >= c.ExpireTime {
			Manager.UnRegister <- c
			break
		}
	}
}

// Read 读取客户端发送过来的消息
func (c *Client) Read() {
	// 出现故障后把当前客户端注销
	defer func() {
		_ = c.Socket.Close()
		Manager.UnRegister <- c
	}()
	for {
		_, data, err := c.Socket.ReadMessage()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		var msg WsMessage
		db := common.GetGorm()
		err = json.Unmarshal(data, &msg)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		switch msg.Type {
		case 6:
			// 如果是心跳监测消息（利用心跳监测来判断对应客户端是否在线）
			resp, _ := json.Marshal(&WsMessage{Type: 6, Data: "pong"})
			c.Start = time.Now() // 重新刷新时间
			c.Send <- resp
		case 1:
			// 获取在线人数
			count := len(Manager.Clients)
			resp, _ := json.Marshal(&WsMessage{Type: 1, Data: count})
			c.Send <- resp
		case 2:
			// 获取消息历史记录(12条最新)
			var chatList []common.TChatRecord
			r1 := db.Model(&common.TChatRecord{}).Order("create_time ASC").Limit(12).Find(&chatList)
			if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
				logger.Error(r1.Error.Error())
				c.Send <- []byte("系统异常")
				return
			}
			// 按照时间正序排序
			sort.Slice(chatList, func(i, j int) bool {
				return chatList[i].CreateTime.Unix() < chatList[j].CreateTime.Unix()
			})
			_data := make(map[string]interface{})
			_data["chatRecordList"] = chatList
			_data["ipAddress"] = c.IpAddress
			_data["ipSource"] = c.IpSource
			resp, _ := json.Marshal(&WsMessage{Type: 2, Data: _data})
			c.Send <- resp
		case 3:
			// 发送文本消息
			var chat common.TChatRecord
			msgByte, _ := json.Marshal(&msg.Data)
			_ = json.Unmarshal(msgByte, &chat)
			r1 := db.Create(&chat)
			if r1.Error != nil {
				logger.Error(r1.Error.Error())
				c.Send <- []byte("系统异常")
				return
			}
			resp, _ := json.Marshal(&WsMessage{Type: 3, Data: chat})
			Manager.Broadcast <- resp

		case 4:
			// 撤回消息
			var req struct {
				ID      int64 `json:"id"`
				IsVoice bool  `json:"isVoice"`
			}
			msgByte, _ := json.Marshal(&msg.Data)
			_ = json.Unmarshal(msgByte, &req)
			switch req.IsVoice {
			case false:
				db.Where("id = ?", req.ID).Delete(&common.TChatRecord{})
			case true:
				var record common.TChatRecord
				r1 := db.Where("id = ?", req.ID).First(&record)
				if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
					logger.Error(r1.Error.Error())
					c.Send <- []byte("系统异常")
					return
				}
				voiceUrl := strings.Split(record.Content, "voice/")
				if len(voiceUrl) == 2 {
					voicePath := common.Conf.App.VoiceDir + voiceUrl[1]
					err := os.Remove(voicePath)
					if err != nil {
						logger.Error(err.Error())
						c.Send <- []byte("系统异常")
						return
					}
				}
				db.Where("id = ?", req.ID).Delete(&common.TChatRecord{})

			}
			c.Send <- data
		case 5:
			// 语音消息
			fmt.Println(msg.Data)

		}

	}
}

// Write 把对应消息写回客户端
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
		Manager.UnRegister <- c
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				// 没有消息则发送空响应
				err := c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Error(err.Error())
					return
				}
				return
			}
			err := c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}
	}
}

var once sync.Once

func GetWsServerManager() *ClientManager {
	once.Do(func() {
		Manager = &ClientManager{
			Broadcast:  make(chan []byte),
			Clients:    make(map[string]*Client),
			Register:   make(chan *Client),
			UnRegister: make(chan *Client), //一定要初始化channel 否则会一直阻塞
		}
	})
	return Manager
}

func (mw *MyWebSocket) WebSocketHandle(ctx *gin.Context) {
	conn, err := (&websocket.Upgrader{
		// 决解跨域问题
		CheckOrigin: func(r *http.Request) bool { return true },
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		logger.Error(err.Error())
		return
	}
	_session, _ := Store.Get(ctx.Request, "CurUser")
	userid := _session.Values["a_userid"]
	ip := ctx.ClientIP()
	addr, err := common.GetIpAddressAndSource(ip)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		logger.Error(err.Error())
		return
	}
	ua := ctx.GetHeader("User-Agent")
	id := ip + ua
	idMd5 := fmt.Sprintf("%x", md5.Sum([]byte(id)))
	if _, ok := Manager.Clients[idMd5]; !ok {
		client := &Client{
			ID:     idMd5,
			Socket: conn, Send: make(chan []byte),
			IpAddress:  ip,
			IpSource:   addr.Data.Province,
			UserId:     userid,
			Start:      time.Now(),
			ExpireTime: time.Minute * 1,
		}
		Manager.Register <- client
		go client.Read()
		go client.Write()
		go client.Check()
	}
}

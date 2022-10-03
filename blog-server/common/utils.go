package common

import (
	"errors"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

/*
* @author:AxisZql
* @date: 2022-3-16 11:45 PM
* @desc: 工具模块
 */

//========= 密码工具

// EncryptionPwd hash加密密码
func EncryptionPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		logger.Debug(fmt.Sprintf("加密密码错误:%v", err))
		return "", err
	}
	return string(hash), nil
}

// VerifyPwd 验证密码
func VerifyPwd(hash, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if err != nil {
		logger.Debug(fmt.Sprintf("密码验证错误:%v", err))
		return false
	}
	return true
}

// IpInfo ========= IP工具
type IpInfo struct {
	Code int64 `json:"code"`
	Data struct {
		Isp      string `json:"isp"` //运营商
		Province string `json:"province"`
		City     string `json:"city"`
		Detail   string `json:"detail"`
	} `json:"data"`
}

// GetIpAddressAndSource 获取IP地址和地区
func GetIpAddressAndSource(ip string) (*IpInfo, error) {
	var _ipInfo IpInfo
	abs, _ := filepath.Abs(".")
	path := filepath.Join(abs, "axisIp.xdb")
	region, err := xdb.NewWithFileOnly(path)
	defer func() {
		region.Close()
	}()
	if err != nil {
		return nil, err
	}
	locate, err := region.SearchByStr(ip)
	if err != nil {
		return nil, err
	}
	addr := strings.Split(locate, "|")
	if len(addr) == 0 {
		return nil, errors.New("ERROR")
	}
	_ipInfo.Data.Province = addr[0]
	_ipInfo.Data.City = addr[0]
	if len(addr) > 1 {
		_ipInfo.Data.Detail = addr[1]
	}
	return &_ipInfo, nil
}

//======转换菜单数据为前端接受的格式

type MenuPart struct {
	Name      interface{} `json:"name"`
	Path      string      `json:"path"`
	Component interface{} `json:"component"`
	Icon      interface{} `json:"icon"`
	Hidden    interface{} `json:"hidden"`
	Children  []MenuPart  `json:"children"`
}

func findMenuChild(data []VUserMenu, parentId int64) ([]MenuPart, []VUserMenu) {
	if len(data) == 0 {
		return nil, data
	}
	h := map[int]bool{0: false, 1: true}
	list := make([]MenuPart, 0)
	for i := 0; i < len(data); i++ {
		val := data[i]
		if val.ParentId == parentId {
			r := data[i+1:]
			data = data[:i]
			data = append(data, r...)
			i-- //取出一个i不能增加
			child, _ := findMenuChild(data, val.MenuId)
			t := MenuPart{
				Name:      val.Name,
				Path:      val.Path,
				Component: val.Component,
				Icon:      val.Icon,
				Hidden:    h[val.IsHidden],
				Children:  child,
			}
			list = append(list, t)
		}
	}
	if len(list) == 0 {
		return nil, data
	}
	return list, data
}

func ConvertMenuType(data []VUserMenu) []MenuPart {
	h := map[int]bool{0: false, 1: true}
	mList := make([]MenuPart, 0)
	for len(data) != 0 {
		for i := 0; i < len(data); i++ {
			val := data[i]
			if val.Component != "Layout" && val.ParentId <= 0 {
				t := MenuPart{
					Path:      val.Path,
					Component: "Layout",
					Hidden:    h[val.IsHidden],
					Children: []MenuPart{
						{
							Name:      val.Name,
							Component: val.Component,
							Icon:      val.Icon,
							Children:  nil,
						},
					},
				}
				mList = append(mList, t)
				r := data[i+1:]
				data = data[:i]
				data = append(data, r...)
				continue
			}
			if val.ParentId <= 0 {
				//先删除出出表节点
				r := data[i+1:]
				data = data[:i]
				data = append(data, r...)
				child, cdata := findMenuChild(data, val.MenuId)
				data = cdata
				t := MenuPart{
					Name:      val.Name,
					Path:      val.Path,
					Component: val.Component,
					Icon:      val.Icon,
					Hidden:    h[val.IsHidden],
					Children:  child,
				}
				mList = append(mList, t)
			}

		}
	}

	return mList
}

//======= 转换评论数据为前端接受的格式

type replyDTO struct {
	ID             int64     `json:"id"`
	ParentId       int64     `json:"parentId"`
	UserId         int64     `json:"userId"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	WebSite        string    `json:"webSite"`
	ReplyUserId    int64     `json:"replyUserId"`
	ReplyNickname  string    `json:"replyNickname"`
	ReplyWebSite   string    `json:"replyWebSite"`
	CommentContent string    `json:"commentContent"`
	LikeCount      int64     `json:"likeCount"`
	CreateTime     time.Time `json:"createTime"`
}

type doneComment struct {
	ID             int64      `json:"id"`
	UserId         int64      `json:"userId"`
	Nickname       string     `json:"nickname"`
	Avatar         string     `json:"avatar"`
	WebSite        string     `json:"webSite"`
	CommentContent string     `json:"commentContent"`
	LikeCount      int64      `json:"likeCount"`
	CreateTime     time.Time  `json:"createTime"`
	ReplyCount     int64      `json:"replyCount"`
	ReplyDTOList   []replyDTO `json:"replyDTOList"`
}

type DoneCommentAddCount struct {
	RecordList []doneComment `json:"recordList"`
	Count      int64         `json:"count"`
}

func GetCommentLikeCountById(id int64) (int64, error) {
	db := GetGorm()
	var likeCount int64
	r1 := db.Model(&TLike{}).Where("object = ? AND like_id = ?", "t_comment", id).Count(&likeCount)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		return 0, r1.Error
	}
	return likeCount, nil
}

//=========发送注册验证码邮件

// GetRandStr 生成n位随机字符串验证码
func GetRandStr(n int) (code string) {
	chars := `ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789`
	charsLen := len(chars)
	if n > 6 {
		n = 6
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		code += chars[randIndex : randIndex+1]
	}
	return code
}

func SentCodeByEmail(code string, toUserEmail ...string) error {
	mailTo := make([]string, 0) //收件人列表
	mailTo = append(mailTo, toUserEmail...)
	title := `AXIS-BLOG注册验证码`
	body := fmt.Sprintf(`Hi👋,您的验证码为:「 <a>%v</a> 」,验证码有效时间为5分钟,请不要将验证码告诉他人喔😉`, code)

	m := gomail.NewMessage()
	m.SetHeader(`From`, Conf.Mail.Username)
	m.SetHeader(`To`, mailTo...)
	m.SetHeader(`Subject`, title)
	m.SetBody(`text/html`, body)
	err := gomail.NewDialer(Conf.Mail.Host, Conf.Mail.Port, Conf.Mail.Username, Conf.Mail.Password).DialAndSend(m)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

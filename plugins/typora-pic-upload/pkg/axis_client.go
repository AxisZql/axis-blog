package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var once sync.Once
var once2 sync.Once

type Cookie struct {
	Value []*http.Cookie
	mutex sync.RWMutex
}

var (
	// host    = "http://127.0.0.1:9080"
	host    = "you server api address"
	Client  *http.Client
	Client2 *http.Client
	cookie  Cookie
)

// 存储图片目标服务器的客户端
func defaultClient() *http.Client {
	once2.Do(func() {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		jar, _ := cookiejar.New(nil)
		Client = &http.Client{
			Timeout:   time.Minute * 2,
			Transport: transport,
			Jar:       jar,
		}
	})
	return Client
}

// 请求外链图片的的客户端
func defaultClient2() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		jar, _ := cookiejar.New(nil)
		Client2 = &http.Client{
			Timeout:   time.Minute * 2,
			Transport: transport,
			Jar:       jar,
		}
	})
	return Client2
}

// 登陆请求
func doRequestLogin(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	client := defaultClient()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := client.Do(req)
	cookie.mutex.Lock()
	cookie.Value = resp.Cookies()
	cookie.mutex.Unlock()
	return resp, err
}

// 获取外链图片请求
func doRequestOtherPic(method, url string, header http.Header) (*http.Response, error) {
	client := defaultClient2()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := client.Do(req)
	return resp, err
}

func Login() error {
	postValue := url.Values{
		"username": {"you blog user"},
		"password": {"you blog user password"},
	}
	postString := postValue.Encode()

	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}

	req := bytes.NewBuffer([]byte(postString))
	resp, err := doRequestLogin("POST", host+"/login", header, req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	mp := make(map[string]interface{}, 0)
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &mp)
	if err != nil {
		return err
	}
	return nil
}

// UpLoadPic
//ty: 上传文件对象为链接还是文件路径
//file: 文件链接或者路径名
func UpLoadPic(ty string, file string) (*http.Response, error) {
	if ty == "url" {
		res, err := doUpLoadOtherReq(file)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if ty == "file" {
		res, err := doUpLoadLocalPicReq(file)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, errors.New("上传图片来源暂时不支持")

}

func doUpLoadOtherReq(fileUrl string) (*http.Response, error) {
	cookie.mutex.RLock()
	defer cookie.mutex.RUnlock()
	if len(cookie.Value) == 0 {
		return nil, errors.New("之前登陆获取的cookie不能为空")
	}
	cookieValue := cookie.Value[0].Value

	header := http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36"},
	}
	data, err := doRequestOtherPic("GET", fileUrl, header)
	if err != nil {
		err = errors.Wrap(err, "获取外链图片失败")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", "pic.png")
	if err != nil {
		return nil, err
	}
	io.Copy(part, data.Body)

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	client := defaultClient()

	_url := host + `/admin/articles/images`
	req, _ := http.NewRequest("POST", _url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	cookie1 := &http.Cookie{Name: "ticket", Value: cookieValue}

	req.AddCookie(cookie1)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 上传本地文件
func doUpLoadLocalPicReq(path string) (*http.Response, error) {
	cookie.mutex.RLock()
	defer cookie.mutex.RUnlock()
	if len(cookie.Value) == 0 {
		return nil, errors.New("之前登陆获取的cookie不能为空")
	}
	cookieValue := cookie.Value[0].Value

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	client := defaultClient()

	_url := host + `/admin/articles/images`
	req, _ := http.NewRequest("POST", _url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	cookie1 := &http.Cookie{Name: "ticket", Value: cookieValue}

	req.AddCookie(cookie1)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

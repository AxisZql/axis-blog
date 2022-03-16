package common

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
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

//========= IP工具
type ipInfo struct {
	Code int64 `json:"code"`
	Data struct {
		Isp      string `json:"isp"` //运营商
		Province string `json:"province"`
		City     string `json:"city"`
	} `json:"data"`
}

// GetIpAddressAndSource 获取IP地址和地区
func GetIpAddressAndSource(ip string) (*ipInfo, error) {
	var header = http.Header{
		"Authorization": []string{fmt.Sprintf("APPCODE %s", Conf.Ip.AppCode)},
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf(`http://cz88.rtbasia.com/search?ip=%s`, ip)
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户地理位置失败:%v", err))
		return nil, err
	}
	defer resp.Body.Close()
	_ipInfo := ipInfo{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &_ipInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("解析IP信息失败:%v", err))
		return nil, err
	}
	return &_ipInfo, nil
}

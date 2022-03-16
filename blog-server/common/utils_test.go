package common

import (
	"fmt"
	"testing"
)

// 测试通过IP获取地理位置的工具
func TestGetLocationByIp(t *testing.T) {
	ipInfo, _ := GetIpAddressAndSource("183.42.38.112")
	fmt.Println(ipInfo)
	if ipInfo == nil || ipInfo.Data.Province == "" {
		t.Fatal("通过IP获取地理位置失败")
	}

}

func TestEncryptionPwd(t *testing.T) {
	s, err := EncryptionPwd("1234567")
	fmt.Println(s)
	if err != nil {
		t.Fatal("加密密码失败")
	}
}

func TestVerifyPwd(t *testing.T) {
	ok := VerifyPwd("$2a$04$fVuFvZfMP6N6joKZw0Is1ulyZOEVgsZxtKWrd8jaL5xzJQ6nkzNT.", "1234567")
	if !ok {
		t.Fatal("验证密码失败")
	}
}

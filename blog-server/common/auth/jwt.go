package auth

import (
	"blog-server/common"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/*
@author:AxisZql
@date:2022-3-26
@desc:jwt工具
*/

var jwtSecret []byte

func init() {
	jwtSecret = []byte(common.Conf.Jwt.Key)
}

// Claims （记录用户实体）
type Claims struct {
	AUserID int64 `json:"a_userid"`
	Role    int64 `json:"role"`
	jwt.StandardClaims
}

// JwtEnc 生成token
func JwtEnc(aUserid int64, role int64) (string, error) {
	now := time.Now()
	expireTime := now.Add(2 * time.Hour) //过期时间为2.0h
	claims := Claims{
		AUserID: aUserid,
		Role:    role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "axiszql.com",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//用密钥进行签名
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

//JwtDec JWT解密
func JwtDec(token string) (*Claims, error) {
	// 进行解析鉴权声明
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

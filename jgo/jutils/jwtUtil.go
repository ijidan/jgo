package jutils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//创建TOKEN
func CreateToken(uid int64, secret string) (string, error) {
	claim := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(3 * 24 * time.Hour).Unix(),
	}
	nt := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return nt.SignedString([]byte(secret))
}

//刷新TOKEN
func RefreshToken(tokenStr string, secret string) (string, error) {
	claim, err := VerifyToken(tokenStr)
	if claim == nil || err != nil {
		err = errors.New("parse token error")
		return "", err
	}
	uid := claim["uid"]
	return CreateToken(uid.(int64), secret)
}

//解析TOKEN
func ParseToken(tokenStr string, secret string) (map[string]interface{}, error) {
	getSecret := func(_secret string) func(token *jwt.Token) (i interface{}, err error) {
		return func(token *jwt.Token) (i interface{}, err error) {
			return []byte(_secret), nil
		}
	}(secret)
	token, err := jwt.Parse(tokenStr, getSecret)
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapClaim")
		return nil, err
	}
	if !token.Valid {
		err = errors.New("token is invalid")
		return nil, err
	}
	return claim, nil
}

const secret = "__jd__"

//生成TOKEN
func GenToken(uid int64) (string, error) {
	return CreateToken(uid, secret)
}

//校验TOKEN
func VerifyToken(tokenStr string) (map[string]interface{}, error) {
	return ParseToken(tokenStr, secret)
}

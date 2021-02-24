package router

import (
	"errors"
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jgo/jgo/jutils"
	"github.com/ijidan/jnet/jnet"
	"log"
	"strings"
)

//日志中间件
func LogMiddleware(c *jcontext.Context) (next bool, err error) {
	log.Println("日志中间件....")
	return true, nil
}

//同步cookie中间件
func BridgeAccountCookieSyncMiddleware(c *jcontext.Context) (next bool, err error) {
	log.Println("同步cookie中间件...")
	request := jnet.NewRequest()
	request.SetGCookies(c.Request.Cookies())
	return true, nil
}

//校验token中间件
func verifyTokenMiddleware(c *jcontext.Context)(next bool,err error)  {
	log.Println("校验token中间件...")
	auth:=c.Request.Header.Get("Authorization")
	token:=""
	keyStr:="Bearer "
	if len(auth)>0 && strings.Contains(auth,keyStr){
		token=strings.ReplaceAll(auth,keyStr,"")
	}
	if len(auth)==0{
		return false,errors.New("token不存在")
	}
	claim, _err :=jutils.VerifyToken(token)
	if claim==nil|| _err!=nil{
		return false,errors.New("token解析失败")
	}
	return true,nil
}

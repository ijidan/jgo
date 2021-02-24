package controller

import (
	"github.com/ijidan/jgo/business/account"
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jnet/jnet"
	"net/http"
	"time"
)

//用户相关
type UserController struct {
}

//用户注册
func (c *UserController) Reg(ctx *jcontext.Context) {
	ctx.JsonSuccess("用户注册", nil, "")
}

//用户登录
func (c *UserController) Login(ctx *jcontext.Context) {
	accountUser := account.User{}
	loginUserId := accountUser.GetLoginUserId()
	if loginUserId > 0 {
		ctx.JsonFail(1, "用户已经登录", nil, "")
	} else {
		accountUser.LoginByAccount("18025473730", "123456")
		//设置cookie
		c.setCookie(ctx.Writer, false)
		ctx.JsonSuccess("用户登录成功", nil, "")
	}

}

//用户信息
func (c *UserController) Info(ctx *jcontext.Context) {
	accountUser := account.User{}
	rsp := accountUser.GetInfo(0, 0, "")
	if rsp.Success() {
		data := rsp.GetData()
		dataMap := data.(map[string]interface{})
		ctx.JsonSuccess("用户信息", dataMap, "")
	} else {
		ctx.JsonFail(1, "用户未登录", nil, "")
	}
}

//退出登录
func (c *UserController) Logout(ctx *jcontext.Context) {
	accountUser := account.User{}
	loginUserId := accountUser.GetLoginUserId()
	if loginUserId > 0 {
		accountUser.Logout()
		c.setCookie(ctx.Writer, true)
		ctx.JsonSuccess("退出登录", nil, "")
	} else {
		ctx.JsonFail(1, "用户未登录", nil, "")
	}

}

//设置cookie
func (c *UserController) setCookie(w http.ResponseWriter, isExpires bool) {
	//设置cookie
	request := jnet.NewRequest()
	gCookies := request.GetGCookies()
	cookieNum := len(gCookies)
	if cookieNum > 0 {
		for idx := 0; idx < cookieNum; idx++ {
			currCookie := gCookies[idx]
			if isExpires == true {
				currCookie.Expires = time.Unix(1, 0)
				currCookie.MaxAge = -1
			}
			http.SetCookie(w, gCookies[idx])
		}
	}
}

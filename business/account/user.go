package account

import (
	"github.com/ijidan/jnet/jnet"
	"reflect"
	"strconv"
)

//用户相关
type User struct {
	Common
	jnet.BaseService
}

//账号密码登录
func (u *User) LoginByAccount(account string, password string) jnet.Response {
	path := "/user/index/login"
	url := u.buildReqUrl(path, "")
	param := map[string]string{"account": account, "password": password}
	reqParam := u.handleParam(param)
	u.SetRequest()
	u.SetUseProxy(true)
	response := u.SendPostRequest(url, reqParam)
	rsp := u.computeResponse(response)
	return rsp
}

//获取用户信息
func (u *User) GetInfo(HFUid int64, uid int64, fields string) jnet.Response {
	path := "/user/index/getInfo"
	url := u.buildReqUrl(path, "")
	param := make(map[string]string)
	if HFUid > 0 {
		param["hf_uid"] = strconv.FormatInt(HFUid, 10)
	}
	if uid > 0 {
		param["uid"] = strconv.FormatInt(uid, 10)
	}
	if len(fields) > 0 {
		param["fields"] = fields
	}
	reqParam := u.handleParam(param)
	u.SetRequest()
	u.SetUseProxy(true)
	response := u.SendPostRequest(url, reqParam)
	rsp := u.computeResponse(response)
	return rsp
}

//获取登录用户ID
func (u *User) GetLoginUserId() int64 {
	rsp := u.GetInfo(0, 0, "")
	if rsp.Fail() {
		return 0
	}
	data := rsp.GetData()
	k := reflect.TypeOf(data).Kind()
	if k != reflect.Map {
		return 0
	}
	dataMap := data.(map[string]interface{})
	uid := dataMap["hf_uid"]
	return int64(uid.(float64))

}

//退出登录
func (u *User) Logout() jnet.Response {
	path := "/user/index/logout"
	url := u.buildReqUrl(path, "")
	param := make(map[string]string)
	reqParam := u.handleParam(param)
	u.SetRequest()
	u.SetUseProxy(true)
	response := u.SendPostRequest(url, reqParam)
	rsp := u.computeResponse(response)
	return rsp
}

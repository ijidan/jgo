package account

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ijidan/jgo/jgo/jutils"
	"github.com/ijidan/jnet/jnet"
	"strconv"
	"time"
)

const signSalt = ""
const domainUser = "" //域名
const businessHf = ""

type Common struct {
	request *jnet.Request
}

//构造请求URL
func (c *Common) buildReqUrl(path string, domain string) string {
	if domain == "" {
		domain = domainUser
	}
	return domain + path
}

//处理参数
func (c *Common) handleParam(param map[string]string) map[string]interface{} {
	timeStamp := time.Now().Unix()
	timeStr := strconv.FormatInt(timeStamp, 10)
	commonParam := map[string]string{"business": businessHf, "time": timeStr}
	mergedParam := commonParam
	if len(param) > 0 {
		for k, v := range param {
			mergedParam[k] = v
		}
	} else {
		mergedParam = commonParam
	}
	sign := c.computeSign(mergedParam)
	mergedParam["sign"] = sign

	//数据格式转化
	convertedParam := make(map[string]interface{})
	for k, v := range mergedParam {
		convertedParam[k] = v
	}
	return convertedParam
}

//计算签名
//计算签名
func (c *Common) computeSign(param map[string]string) string {
	sortedParam := jutils.KSort(param)
	str := sortedParam + signSalt
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

//设置请求参数
func (c *Common) SetRequest() {
	request := jnet.NewRequest()
	c.request = request
}

//获取域名
func (c *Common) GetDomain() string {
	return domainUser
}

//响应
func (c *Common) computeResponse(rsp jnet.Response) jnet.Response {
	isFail := rsp.Fail()
	if isFail {
		return jnet.NewResponse(rsp.GetCode(), rsp.GetMessage(), rsp.GetData(), rsp.GetPrompt())
	}
	body := rsp.GetData()
	bodyString := body.(string)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bodyString), &m); err != nil {
		return jnet.NewResponse(jnet.JsonParseFail, err.Error(), nil, "")
	}
	r := c.request
	//cookie处理
	r.GCookies = r.GCookieJar.Cookies(r.Req.URL)
	//返回结果处理
	stateInterface := m["state"]
	dataInterface := m["data"]
	msgInterface := m["msg"]
	//数据格式转换
	state := int64(stateInterface.(float64))
	msg := msgInterface.(string)
	if state > 0 {
		return jnet.NewResponse(jnet.ReturnFail, msg, dataInterface, "")
	}
	return jnet.NewResponse(jnet.Success, msg, m["data"], "")
}

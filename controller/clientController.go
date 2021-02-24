package controller

import (
	"encoding/json"
	. "github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jws/jws"
	"github.com/sirupsen/logrus"
)

//默认
type ClientController struct {
}

//首页
func (c *ClientController) GetAll(ctx *Context) {
	tools := jws.Tools{}
	allClientList := tools.GetAllClientList()

	data := make(map[string]interface{})
	data["total"] = len(allClientList)
	data["rows"] = allClientList

	dataJson, err := json.Marshal(data)
	if err != nil {
		errMsg := err.Error()
		logrus.Println(errMsg)
	}
	ctx.JsonRaw(string(dataJson))
}

//发信息
func (c *ClientController) SendMessage(ctx *Context) {
	clientId := ctx.GetPost("client_id", "")
	messageContent := ctx.GetPost("message_content", "")
	if messageContent == "" || clientId == "" {
		ctx.JsonFail(1, "参数错误", nil, "")
		return
	}
	tools := jws.Tools{}
	client := tools.GetClientByClientId(clientId)
	if client == nil {
		ctx.JsonFail(1, "客户端不存在", nil, "")
		return
	}
	tools.SendTextMessageToClient(client, messageContent)
	ctx.JsonSuccess("发送成功", nil, "")

}

//踢下线
func (c *ClientController) KickOff(ctx *Context) {
	clientId := ctx.GetPost("client_id", "")
	if clientId == "" {
		ctx.JsonFail(1, "参数错误", nil, "")
		return
	}
	tools := jws.Tools{}
	client := tools.GetClientByClientId(clientId)
	if client == nil {
		ctx.JsonFail(1, "客户端不存在", nil, "")
		return
	}
	tools.KickOff(client, "服务端主动关闭")
	ctx.JsonSuccess("操作成功", nil, "")

}

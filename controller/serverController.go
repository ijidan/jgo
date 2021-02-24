package controller

import (
	"encoding/json"
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jws/jws"
	"github.com/sirupsen/logrus"
)

//默认
type ServerController struct {
}

//首页
func (c *ServerController) GetAll(ctx *jcontext.Context) {
	allServerList := jws.GetServerList()
	data := make(map[string]interface{})
	data["total"] = len(allServerList)
	data["rows"] = allServerList

	dataJson, err := json.Marshal(data)
	if err != nil {
		errMsg := err.Error()
		logrus.Println(errMsg)
	}
	ctx.JsonRaw(string(dataJson))
}

//批量发送信息
func (c *ServerController) BatchSendMessage(ctx *jcontext.Context) {
	serverId := ctx.GetPost("server_id", "")
	messageContent := ctx.GetPost("message_content", "")
	if messageContent == "" || serverId == "" {
		ctx.JsonFail(1, "参数错误", nil, "")
		return
	}
	server := jws.GetServerByServerId(serverId)
	if server == nil {
		ctx.JsonFail(1, "服务器不存在", nil, "")
		return
	}
	jws.SendTextMessageToServer(server, messageContent)
	ctx.JsonSuccess("发送成功", nil, "")
}

//关闭服务器
func (c *ServerController) CloseServer(ctx *jcontext.Context) {
	serverId := ctx.GetPost("server_id", "")
	if serverId == "" {
		ctx.JsonFail(1, "参数错误", nil, "")
		return
	}
	server := jws.GetServerByServerId(serverId)
	if server == nil {
		ctx.JsonFail(1, "服务器不存在", nil, "")
		return
	}
	server.Close()
	ctx.JsonSuccess("关闭成功", nil, "")
}

package controller

import (
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jproto/jproto/call"
)

//默认
type IndexController struct {
}

//首页
func (c *IndexController) Index(ctx *jcontext.Context) {
	ctx.RenderNoLayoutTemplate("index/index.html", nil)
}

//proto调用
func (c *IndexController) Call(ctx *jcontext.Context) {
	//span, _ := opentracing.StartSpanFromContext(ctx.TracingCtx, "jaeger-span-call")
	//req := map[string]string{"name": param}
	//replay := mapData


	param := "jidan"
	hello := call.HelloCall{}
	data, _ := hello.SayHello(param)
	mapData := make(map[string]interface{})
	mapData["result"] = data

	//span.SetTag("request", "req:name:jidan")
	//span.SetTag("reply", "rep:hello jidan")
	//span.Finish()

	//返回结果
	ctx.JsonSuccess("", mapData, "")
}

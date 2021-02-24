package controller

import (
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jgo/jgo/jutils"
)

//默认
type SystemController struct {
}


//首页
func (c *SystemController) Error404(ctx *jcontext.Context) {
	timeStr:=(jutils.NewTimeUtil()).GetCurrentTime()
	data:=map[string]string{"message":"404","time":timeStr}
	ctx.RenderNoLayoutTemplate("system/error.html", data)
}

//内部错误
func (c *SystemController) ErrorInner(ctx *jcontext.Context) {
	timeStr:=(jutils.NewTimeUtil()).GetCurrentTime()
	data:=map[string]string{"message":"inner error","time":timeStr}
	ctx.RenderNoLayoutTemplate("system/error.html", data)
}


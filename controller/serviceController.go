package controller

import (
	"github.com/ijidan/jgo/jgo/jcontext"
)

//服务相关
type ServiceController struct {
}

//consul check
func (c *ServiceController) ConsulCheck(ctx *jcontext.Context) {
	ctx.Text("consul check success")
}

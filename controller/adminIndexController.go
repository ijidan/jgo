package controller

import "github.com/ijidan/jgo/jgo/jcontext"

//默认
type AdminIndexController struct {
}

//首页
func (c *AdminIndexController) Index(ctx *jcontext.Context) {
	ctx.RenderL2("admin/index.html", nil)
}

//客户端列表
func (c *AdminIndexController) ClientList(ctx *jcontext.Context) {
	ctx.RenderL1("admin/client_list.html", nil)
}

//服务器列表
func (c *AdminIndexController) ServerList(ctx *jcontext.Context) {
	ctx.RenderL1("admin/server_list.html", nil)
}

//默认首页
func (c *AdminIndexController) DefaultIdx(ctx *jcontext.Context) {
	ctx.RenderL1("default_index.html", nil)
}

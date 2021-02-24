package controller

import (
	"github.com/ijidan/jgo/jgo/jcontext"
)

//聊天相关
type ChatController struct {
}

//聊天首页
func (c *ChatController) Index(ctx *jcontext.Context) {
	ctx.RenderL1("chat/index.html", nil)
}

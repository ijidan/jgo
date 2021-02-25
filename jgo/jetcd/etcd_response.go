package jetcd

import "github.com/ijidan/jnet/jnet"

//响应类
type Response struct {
	responseString string
	code           int64
	message        string
	prompt         string
	data           interface{}
	cookies        []string
	header         []string
}

//构造函数
func (r *Response) Construct(code int64, message string, data interface{}, prompt string) {
	r.code = code
	r.message = message
	r.data = data
	r.prompt = prompt
}

//是否成功
func (r *Response) Success() bool {
	return r.code == jnet.Success
}

//是否失败
func (r *Response) Fail() bool {
	return r.code != jnet.Success
}

//设置数据
func (r *Response) SetData(data interface{}) {
	r.data = data
}

//设置cookie
func (r *Response) SetCookies(cookies []string) {
	r.cookies = cookies
}

//设置header
func (r *Response) SetHeader(header []string) {
	r.header = header
}

//获取code
func (r *Response) GetCode() int64 {
	return r.code
}

//获取信息
func (r *Response) GetMessage() string {
	return r.message
}

//获取提示
func (r *Response) GetPrompt() string {
	return r.prompt
}

//获取数据
func (r *Response) GetData() interface{} {
	return r.data
}

//获取cookies
func (r *Response) GetCookies() []string {
	return r.cookies
}

//获取header
func (r *Response) GetHeader() []string {
	return r.header
}

//字符串
func (r *Response) ToString() string {
	return r.responseString
}

//构造响应
func BuildResponse(code int64, message string, data interface{}, prompt string) Response {
	rsp := Response{}
	rsp.Construct(code, message, data, prompt)
	return rsp
}

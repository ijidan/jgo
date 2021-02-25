package jcontext

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//上下文
type Context struct {
	TracingCtx context.Context

	Writer  http.ResponseWriter
	Request *http.Request

	Path       string
	Method     string
	StatusCode int64
}

//获取POST数据
func (c *Context) GetPost(key string, defaultValue string) string {
	_ = c.Request.ParseForm()
	value := ""
	postParams := c.Request.PostForm
	for k, v := range postParams {
		if k == key {
			value = v[0]
			break
		}
	}
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

//获取GET数据
func (c *Context) GetQuery(key string, defaultValue string) string {
	query := c.Request.URL.Query()
	value := query.Get(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

//设置状态码
func (c *Context) SetStatusCode(statusCode int64) {
	c.StatusCode = statusCode
	c.Writer.WriteHeader(int(c.StatusCode))
}

//设置header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//获取模板内容
func (c *Context) GetTemplate(templateName string) string {
	templatePath := c.getTemplatePath(templateName)
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return ""
	}
	return string(content)
}

//渲染无布局模板
func (c *Context) RenderNoLayoutTemplate(templateName string, data interface{}) {
	templatePath := c.getTemplatePath(templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入响应
	c.SetHeader("Content-Type", "text/html")
	_ = tmpl.Execute(c.Writer, data)
}

//渲染有布局模板
func (c *Context) RenderLayoutTemplate(templateName string, data interface{}, layoutName ...string, ) {
	//计算模板文件
	layoutPathList := make([]string, 0)
	for _, item := range layoutName {
		itemPath := c.getTemplatePath(item)
		layoutPathList = append(layoutPathList, itemPath)
	}
	templatePath := c.getTemplatePath(templateName)
	layoutPathList = append(layoutPathList, templatePath)

	t, err := template.ParseFiles(layoutPathList...)
	if err != nil {
		fmt.Printf("parse files failed, err : %v\n", err)
		return
	}
	// 渲染模板
	// 渲染模板时使用ExecuteTemplate函数，需要制定要被渲染的模板名称
	name := strings.ReplaceAll(layoutName[0], ".html", "")
	log.Println(name)
	c.SetHeader("Content-Type", "text/html")
	_ = t.ExecuteTemplate(c.Writer, name, data)
}

//渲染模板1
func (c *Context) RenderL1(templateName string, data interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	c.RenderLayoutTemplate(templateName, data, "layout1.html", "default_css.html", "default_js.html")

}

//渲染模板2
func (c *Context) RenderL2(templateName string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	year := time.Now().Year()
	data["year"] = year
	c.RenderLayoutTemplate(templateName, data, "layout2.html", "layout1.html", "default_css.html", "default_js.html")
}

//重定向
func (c *Context) Redirect(url string, params []string) {
}

//404
func (c *Context) NotFound() {
}

//成功JSON响应
func (c *Context) JsonSuccess(message string, data map[string]interface{}, jumpUrl string) {
	c.SetHeader("Content-Type", "application/json")
	result := c.json(0, message, data, jumpUrl)
	_, _ = c.Writer.Write([]byte(result))
}

//失败JSON响应
func (c *Context) JsonFail(code int64, message string, data map[string]interface{}, jumpUrl string) {
	c.SetHeader("Content-Type", "application/json")
	result := c.json(code, message, data, jumpUrl)
	_, _ = c.Writer.Write([]byte(result))
}

//输出字符串
func (c *Context) Text(data string) {
	c.SetHeader("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte(data))
}

//原始JSON
func (c *Context) JsonRaw(data string) {
	c.SetHeader("Content-Type", "application/json")
	_, _ = c.Writer.Write([]byte(data))

}

//JSONP
func (c *Context) JsonP(code int64, message string, data map[string]interface{}, jumpUrl string, callback string) {
	c.SetHeader("Content-Type", "text/plain")
	result := c.json(code, message, data, jumpUrl)
	if result == "" {
		_, _ = c.Writer.Write([]byte(result))
	} else {
		if callback == "" {
			callback = "_callback"
		}
		result = fmt.Sprintf("%s(%s)", callback, result)
		_, _ = c.Writer.Write([]byte(result))
	}
}

//成功iframe响应
func (c *Context) IFrameWriterSuccess(message string, data map[string]interface{}, jumpUrl string, callback string) {
	c.SetHeader("Content-Type", "text/plain")
	result := c.iFrameResponse(0, message, data, jumpUrl, callback)
	_, _ = c.Writer.Write([]byte(result))
}

//失败iframe响应
func (c *Context) IFrameResponseFail(code int64, message string, data map[string]interface{}, jumpUrl string, callback string) {
	c.SetHeader("Content-Type", "text/plain")
	result := c.iFrameResponse(code, message, data, jumpUrl, callback)
	_, _ = c.Writer.Write([]byte(result))
}

//输出JSON
func (c *Context) json(code int64, message string, data map[string]interface{}, jumpUrl string) string {
	result := c.buildResponseResult(code, message, data, jumpUrl)
	b, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(b)
}

//iframe响应格式
func (c *Context) iFrameResponse(code int64, message string, data map[string]interface{}, jumpUrl string, callback string) string {
	result := c.json(code, message, data, jumpUrl)
	if result == "" {
		return ""
	}
	if callback == "" {
		callback = "_callback"
	}
	html := fmt.Sprintf(`<!doctype html><html lang="en"><head><meta charset="UTF-8" /><title></title><script>
				var frame = null;
				try {
					frame = window.frameElement;
					if(!frame){
						throw("no frame 1");
					}
				} catch(ex){
					try {
						document.domain = location.host.replace(/^[\w]+\./, \'\');
						frame = window.frameElement;
						if(!frame){
							throw("no frame 2");
						}
					} catch(ex){
						if(window.console){
							console.log("i try twice to cross domain. sorry, i m give up...");
						}
					}
				};
				</script><script>frame.%s(%s);</script></head><body></body></html>`, callback, result)
	return html
}

//构造响应结果
func (c *Context) buildResponseResult(code int64, message string, data map[string]interface{}, jumpUrl string) map[string]interface{} {
	result := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
		"jumpUrl": jumpUrl,
	}
	return result
}

//计算静态资源文件路径
func (c *Context) ComputeStaticFilePath(fileName string) string {
	protectedPath := c.getProtectedDir()
	return protectedPath + fileName
}

//获取模板路径
func (c *Context) getTemplatePath(templateName string) string {
	templateDir := c.getTemplateDir()
	templatePath := templateDir + templateName
	return templatePath
}

//获取图片文件路径
func (c *Context) getImgPath(imgName string) string {
	return c.getStaticFilePath("img", imgName)
}

//获取JS文件路径
func (c *Context) getJsPath(jsName string) string {
	return c.getStaticFilePath("js", jsName)
}

//获取CSS文件路径
func (c *Context) getCssPath(cssName string) string {
	return c.getStaticFilePath("css", cssName)
}

//获取模板目录
func (c *Context) getTemplateDir() string {
	protectedPath := c.getProtectedDir()
	pathSep := c.getPathSep()
	dirPath := protectedPath + "template" + pathSep
	return dirPath
}

//获取静态资源文件路径
func (c *Context) getStaticFilePath(cat string, fileName string) string {
	staticDir := c.getStaticDir()
	pathSep := c.getPathSep()
	return staticDir + cat + pathSep + fileName
}

//获取静态资源目录
func (c *Context) getStaticDir() string {
	protectedPath := c.getProtectedDir()
	pathSep := c.getPathSep()
	dirPath := protectedPath + "static" + pathSep
	return dirPath
}

//获取protected目录
func (c *Context) getProtectedDir() string {
	wd, _ := os.Getwd()
	pathSep := c.getPathSep()
	protectedPath := wd + pathSep
	return protectedPath
}

//获取目录分隔符
func (c *Context) getPathSep() string {
	pathSep := string(os.PathSeparator)
	return pathSep
}

func (c *Context) RecoverHandler() {
	if r := recover(); r != nil {

	}
}

//获取实例
func NewContext(w http.ResponseWriter, req *http.Request, tracingCtx context.Context) *Context {
	return &Context{
		TracingCtx: tracingCtx,
		Writer:     w,
		Request:    req,
		Path:       req.URL.Path,
		Method:     req.Method,
	}
}

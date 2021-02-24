package jrouter

import (
	"net/http"
)

//解析器
type Parser struct {
}

const (
	modelNormal     = iota //普通模式 http://localhost/index.go?c=ctrl&a=action
	modelPathInfo          //PATH INFO 模式 http://localhost/index.go/ctrl/act
	modelRewrite           //Rewrite 模式 http://localhost/ctrl/act
	modelCompatible        //兼容模式 http://localhost/index.go?r=ctrl/act
)

//解析控制器和方法
func (p *Parser) GetMatch(r *http.Request) string {
	match := p.getRewriteModelMatch(r)
	return match
}

//普通模式
func (p *Parser) getNormalModelMatch(r *http.Request) string {
	query := r.URL.Query()
	return query["c"][0] + "/" + query["a"][0]
}

//PathInfo模式
func (p *Parser) getPathInfoModelMatch(r *http.Request) string {
	query := r.URL.Query()
	match := ""
	for k, _ := range query {
		match = k
		break
	}
	return match
}

//rewrite模式
func (p *Parser) getRewriteModelMatch(r *http.Request) string {
	path := r.URL.Path
	return path
}

//兼容模式
func (p *Parser) getCompatibleModelMatch(r *http.Request) string {
	query := r.URL.Query()
	caList := query["r"]
	ca := caList[0]
	return ca
}

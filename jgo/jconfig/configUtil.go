package jconfig

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

//配置工具
type ConfigUtil struct {
}

//获取配置
func (u *ConfigUtil) GetConfig(key string) string {
	contentStr := u.getConfigContent()
	if len(contentStr) == 0 {
		return ""
	}
	value := gjson.Get(contentStr, key)
	return value.String()
}

//批量获取配置
func (u *ConfigUtil) GetConfigList(keyList []string, valueList ...*string) {
	for idx, key := range keyList {
		value := u.GetConfig(key)
		*valueList[idx] = value
	}
}

//获取配置内容
func (u *ConfigUtil) getConfigContent() string {
	wd, _ := os.Getwd()
	pathSep := string(os.PathSeparator)
	confPath := wd + pathSep + "protected" + pathSep + "config" + pathSep
	fileName := confPath + "default.json"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(content)
}

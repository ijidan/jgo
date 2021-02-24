package jutils

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//字符串转其他类型
func ConvertValueFromString(value string, t reflect.Type) interface{} {
	var convertedValue interface{}
	switch t.Kind() {
	case reflect.Bool:
		convertedValue, _ = strconv.ParseBool(value)
		break
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		convertedValue, _ = strconv.ParseInt(value,10,64)
		break
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		convertedValue, _ = strconv.ParseUint(value, 10, 64)
		break
	case reflect.Float32,
		reflect.Float64:
		convertedValue, _ = strconv.ParseFloat(value, 64)
		break
	case reflect.String:
		convertedValue = value
		break
	default:
		convertedValue = value
	}

	return convertedValue

}

//其他类型转化为字符串
func ConvertValue2String(params map[string]interface{}) map[string]string {
	paramList := make(map[string]string)
	if len(params) > 0 {
		for k, v := range params {
			value := ConvertVal2String(v)
			paramList[k] = value
		}
	}
	return paramList
}

//变量转化为字符串
func ConvertVal2String(v interface{}) string {
	var value string
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		_v := reflect.ValueOf(v).Int()
		value = strconv.FormatInt(_v, 10)
		break
	default:
		value = v.(string)
	}
	return value

}

//JSON转MAP
func ConvertJSON2Map(content string) (r map[string]map[string]interface{}, err error) {
	if len(content) == 0 {
		return nil, errors.New("content empty")
	}
	var result map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}
	return result, nil
}

//字符串转slice
func ConvertString2Slice(key string, sep string) []string {
	var keyList []string
	if find := strings.Contains(key, sep); find {
		keyList = strings.Split(key, sep)
	} else {
		keyList = []string{key}
	}
	return keyList
}

//排序
func KSort(params map[string]string) string {
	var dataParams string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//拼接
	for _, k := range keys {
		dataParams = dataParams + k + "=" + params[k] + "&"
	}
	removedDataParams := dataParams[0 : len(dataParams)-1]
	return removedDataParams
}

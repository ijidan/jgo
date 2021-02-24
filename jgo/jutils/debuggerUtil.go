package jutils

import (
	"log"
	"reflect"
)

//打印变量
func PR(v interface{}) {
	t := reflect.TypeOf(v)
	kind := t.Kind()
	log.Println("****************")
	if kind == reflect.Map {
		value := reflect.ValueOf(v)
		keys := value.MapKeys()
		log.Println("****************")
		for _, key := range keys {
			s := value.MapIndex(key)
			log.Printf("%v", s.Interface())
		}
	} else {
		log.Printf("%v", v)
	}
	log.Println("****************")

}

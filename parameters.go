package iodata

import (
	"reflect"
	"strings"
)

func ParseParameters(params map[string]interface{}, dest interface{}) {
	value := reflect.ValueOf(dest).Elem()
	for key, param := range params {
		parts := strings.Split(key, ".")
		v := value
		for _, p := range parts {
			v = v.FieldByName(p)
			if !v.IsValid() {
				continue
			}
		}
		v.Set(reflect.ValueOf(param))
	}
}

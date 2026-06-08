package commons

import (
	"reflect"
	"strings"
)

func StructToUpdateMap(data interface{}) map[string]interface{} {
	updates := make(map[string]interface{})
	v := reflect.ValueOf(data)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			tag := t.Field(i).Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			tagName := strings.Split(tag, ",")[0]
			updates[tagName] = field.Elem().Interface()
		}
	}
	return updates
}

package commons

import (
	"bokee/internal/mods/interfaces"
	"fmt"

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

// CheckOwnership 校验当前用户是否拥有该资源，否则返回"无权操作"错误
func CheckOwnership(record interfaces.Ownable, userId uint) error {
	if record.GetUserID() != userId {
		return fmt.Errorf("无权操作他人资源")
	}
	return nil
}

package account

import (
	"github.com/nareshganesan/services/shared"
	"reflect"
	"strings"
)

// EntityToMap converts account entity map[string]interface object
func EntityToMap(account *Entity) *map[string]interface{} {
	data := make(map[string]interface{})
	fields := reflect.ValueOf(account).Elem()
	fieldType := fields.Type()
	for i := 0; i < fields.NumField(); i++ {
		fieldName := shared.GetSnakeCase(fieldType.Field(i).Name)
		fieldValue := fields.Field(i).Interface()
		if !strings.Contains(strings.ToLower(fieldType.Field(i).Name), "id") {
			if fieldValue != "" {
				data[fieldName] = fieldValue
			}
		}
	}
	return &data
}

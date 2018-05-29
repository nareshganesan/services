package account

import (
	// "fmt"
	// g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"reflect"
	"strings"
)

// EntityToMap converts account entity map[string]interface object
func EntityToMap(account *Entity) *map[string]interface{} {
	// es := g.GetGlobals()
	// l := es.Log
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
	// use it only for debugging
	// l.WithFields(logrus.Fields{
	// 	"data": data,
	// }).Info("account entity to map")
	return &data
}

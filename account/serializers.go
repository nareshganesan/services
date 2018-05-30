package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// g "github.com/nareshganesan/services/globals"
	// "github.com/sirupsen/logrus"
	// "net/http"
	"reflect"
)

// Serializer serializer for account entity
type Serializer struct {
	Ctx *gin.Context
	Entity
}

// ListSerializer serializer for account entity
type ListSerializer struct {
	Ctx      *gin.Context
	Accounts []Entity
}

// Account represents account data to be exposed by the serializer
type Account struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    string   `json:"roles"`
	Title    string   `json:"title"`
	Phone    string   `json:"phone"`
	Location location `json:"location"`
}

// location object of Account entity
type location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// get creates Account from Serializer object
func (s *Serializer) get() Account {
	account := Account{
		ID:    s.ID,
		Email: s.Email,
		Name:  s.Name,
		Title: s.Title,
		Roles: s.Roles,
	}
	return account
}

// get creates Account array from ListSerializer object
func (s *ListSerializer) get() []Account {
	accounts := []Account{}
	for _, acc := range s.Accounts {
		serializer := Serializer{s.Ctx, acc}
		accounts = append(accounts, serializer.get())
	}
	return accounts
}

// Dump ,serializer api to expose account entity as map[string]interface object
func (s *Serializer) Dump(data map[string]interface{}) (map[string]interface{}, error) {
	out, _ := s.toMap()
	for key, val := range data {
		out[key] = val
	}
	return out, nil
}

// Dump ,serializer api to expose account entity as map[string]interface object
func (s *ListSerializer) Dump(data map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	var accounts []map[string]interface{}
	for _, acc := range s.Accounts {
		serializer := Serializer{s.Ctx, acc}
		account, _ := serializer.toMap()
		accounts = append(accounts, account)
	}
	out["accounts"] = accounts
	for key, val := range data {
		out[key] = val
	}
	return out, nil
}

// toMap converts serializer entity to map[string]interface object
// Ref: https://stackoverflow.com/questions/23589564/function-for-converting-a-struct-to-map-in-golang
func (s *Serializer) toMap() (map[string]interface{}, error) {
	tag := "json"
	// es := g.GetGlobals()
	// l := es.Log
	in := s.get()
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// accept only structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Only struct type is supported; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}
	// use it only for debugging
	// l.WithFields(logrus.Fields{
	// 	"data": out,
	// }).Info("account serializer to map")
	return out, nil
}

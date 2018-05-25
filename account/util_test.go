package account

import (
	"reflect"
	"testing"
)

func getAccountEntity() *Entity {
	var account Entity
	account.Username = "username"
	account.Email = "email@email.com"
	account.Password = "password"
	return &account
}

func getAccountMap() *map[string]interface{} {
	expected := make(map[string]interface{})
	expected["username"] = "username"
	expected["email"] = "email@email.com"
	expected["password"] = "password"
	expected["name"] = ""
	expected["verification_token"] = ""
	expected["roles"] = ""
	expected["is_archived"] = false
	expected["is_verified"] = false
	expected["title"] = ""
	return &expected
}

func TestEntityToMapEquality(t *testing.T) {
	account := getAccountEntity()
	expected := getAccountMap()
	actual := EntityToMap(account)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("TestEntityToMapEquality failed!")
		t.Fail()
	}
}

func TestEntityToMapInEquality(t *testing.T) {
	account := getAccountEntity()
	account.Name = "name"
	expected := getAccountMap()
	actual := EntityToMap(account)
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("TestEntityToMapInEquality failed!")
		t.Fail()
	}
}

package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getAccountEntity() *Entity {
	var account Entity
	account.Email = "email@email.com"
	account.Password = "password"
	account.IsArchived = false
	account.IsArchived = false
	return &account
}

func getAccountMap() *map[string]interface{} {
	expected := make(map[string]interface{})
	expected["email"] = "email@email.com"
	expected["password"] = "password"
	expected["is_archived"] = false
	expected["is_verified"] = false
	return &expected
}

func TestEntityToMapEquality(t *testing.T) {
	account := getAccountEntity()
	expected := getAccountMap()
	actual := EntityToMap(account)
	assert.Equal(t, expected, actual)
}

func TestEntityToMapInEquality(t *testing.T) {
	account := getAccountEntity()
	account.Name = "name"
	expected := getAccountMap()
	actual := EntityToMap(account)
	assert.NotEqual(t, expected, actual)
}

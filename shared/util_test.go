package shared

import (
	"github.com/stretchr/testify/assert"
	// "reflect"
	// "fmt"
	"testing"
)

func TestGetHash(t *testing.T) {
	hash, _ := GetHash("password")
	expected := true
	assert.Equal(t, expected, VerifyHash("password", hash))
}

func TestVerifyHash(t *testing.T) {
	hash, _ := GetHash("password")
	assert.Equal(t, true, VerifyHash("password", hash))
}

func TestGetSnakeCase(t *testing.T) {
	data := "apiHandler"
	expected := "api_handler"
	actual := GetSnakeCase(data)
	assert.Equal(t, expected, actual)
}

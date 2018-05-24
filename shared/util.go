package shared

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"unicode"
)

// GetSnakeCase returns snakecase string from camel case string
func GetSnakeCase(name string) string {
	var words []string
	var snakeCase string

	l := 0
	for s := name; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
		snakeCase += strings.ToLower(s[:l]) + "_"
	}
	snakeCase = strings.TrimRight(snakeCase, "_")
	return snakeCase
}

// DefaultInt returns safely extracts param with default int value
func DefaultInt(ctx *gin.Context, param string, def int) int {
	defString := strconv.Itoa(def)
	if val, err := strconv.Atoi(ctx.DefaultQuery(param, defString)); err == nil {
		return val
	}
	return def
}

// GetHash returns hash value for given string
func GetHash(value string) (string, error) {
	fmt.Println("Hashing value")
	hash, err := bcrypt.GenerateFromPassword([]byte(value), 14)
	return string(hash), err
}

// VerifyHash returns true if hash and the string are same
func VerifyHash(value, hash string) bool {
	fmt.Println("verifyHash")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}

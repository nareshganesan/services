package account

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

func TestLoginValidCredentials(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test1@email.com", "password": "Testpassword#123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusOK
	assert.Equal(t, expected, w.Code)

}

func TestLoginInValidCredentials(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test123", "password": "test@123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusOK
	assert.NotEqual(t, expected, w.Code)

}

func TestLoginEmptyEmail(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "", "password": "test@123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusOK
	assert.NotEqual(t, expected, w.Code)

}

func TestLoginEmptyPassword(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test123", "password": ""}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusOK
	assert.NotEqual(t, expected, w.Code)

}

func TestLoginInvalidData(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("\"foo\":\"bar\", \"bar\":\"foo\"}"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusBadRequest
	assert.Equal(t, expected, w.Code)

}

func TestSignupValidCredentials(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test2@email.com", "password": "Testpassword#123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusOK
	assert.Equal(t, expected, w.Code)

}

func TestSignupExistingUser(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test1@email.com", "password": "Testpassword#123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusUnauthorized
	assert.Equal(t, expected, w.Code)

}

func TestSignupEmptyemail(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "", "password": "test@123"}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusUnauthorized
	assert.Equal(t, expected, w.Code)

}

func TestSignupEmptyPassword(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	data := map[string]string{"email": "test", "password": ""}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusUnauthorized
	assert.Equal(t, expected, w.Code)

}

func TestSignupInvalidData(t *testing.T) {

	// Create a response recorder
	w := httptest.NewRecorder()
	r := GetRouter(true)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString("\"foo\":\"bar\", \"bar\":\"foo\"}"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r.ServeHTTP(w, req)
	expected := http.StatusBadRequest
	assert.Equal(t, expected, w.Code)

}

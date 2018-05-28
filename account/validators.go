package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nareshganesan/services/shared"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// VRegisterRequest validates signup request
func VRegisterRequest(ctx *gin.Context, obj interface{}) (*shared.Response, url.Values) {
	validErrs := url.Values{}
	if invResp, err := shared.ValidateRequest(ctx, obj); err != nil {
		validErrs.Add("error", err.Error())
		return invResp, validErrs
	}
	user := obj.(*Entity)

	if user.Email == "" {
		validErrs.Add("email", "email is required")
	}
	if user.Password == "" {
		validErrs.Add("password", "password is required")
	}
	if user.Password != "" {
		if len(user.Password) < 8 {
			validErrs.Add("password", "minimum password length 8")
		}
		if len(user.Password) > 255 {
			validErrs.Add("password", "maximum password length 255")
		}
		// following are the basic password rules
		// 1. should contain atleast one special character
		passwordSpecialCharRule := "#@,_-~!*^?+=:;"
		if !strings.ContainsAny(user.Password, passwordSpecialCharRule) {
			mess := fmt.Sprintf("atleast one of the following special characters %s should be present in password", passwordSpecialCharRule)
			validErrs.Add("password", mess)
		}
		// 2. should contain atleast one number
		passwordNumberRule := regexp.MustCompile("[0-9]+")
		if len(passwordNumberRule.FindAllString(user.Password, 1)) == 0 {
			validErrs.Add("password", "password should contain atleast one number")
		}
		// 3. should contain atleast one small alphabet
		passwordSmallCharRule := regexp.MustCompile("[a-z]+")
		if len(passwordSmallCharRule.FindAllString(user.Password, 1)) == 0 {
			validErrs.Add("password", "password should contain atleast one small alphabet [a-z]")
		}
		// 4. should contain atleast one capital alphabet
		passwordCapitalCharRule := regexp.MustCompile("[A-Z]+")
		if len(passwordCapitalCharRule.FindAllString(user.Password, 1)) == 0 {
			validErrs.Add("password", "password should contain atleast one capital alphabet [A-Z]")
		}
	}
	data := make(map[string]interface{})
	if len(validErrs) > 0 {
		data["error"] = validErrs
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
		data["status"] = http.StatusText(http.StatusUnauthorized)
		invResp := &shared.Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, validErrs
	}
	return nil, nil
}

// VLoginRequest validates login request
func VLoginRequest(ctx *gin.Context, obj interface{}) (*shared.Response, url.Values) {
	validErrs := url.Values{}
	if invResp, err := shared.ValidateRequest(ctx, obj); err != nil {
		validErrs.Add("error", err.Error())
		return invResp, validErrs
	}
	user := obj.(*Entity)

	if user.Email == "" {
		validErrs.Add("email", "email is required")
	}
	if user.Password == "" {
		validErrs.Add("password", "password is required")
	}
	data := make(map[string]interface{})
	if len(validErrs) > 0 {
		data["error"] = validErrs
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
		data["status"] = http.StatusText(http.StatusUnauthorized)
		invResp := &shared.Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, validErrs
	}
	return nil, nil
}

// VUpdateAccountRequest validates update account request
func VUpdateAccountRequest(ctx *gin.Context, obj interface{}) (*shared.Response, url.Values) {
	validErrs := url.Values{}
	if invResp, err := shared.ValidateRequest(ctx, obj); err != nil {
		validErrs.Add("error", err.Error())
		return invResp, validErrs
	}
	user := obj.(*Entity)

	if (user.Username == "") && (user.Email == "") {
		validErrs.Add("email", "email is required")
	}
	data := make(map[string]interface{})
	if len(validErrs) > 0 {
		data["error"] = validErrs
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
		data["status"] = http.StatusText(http.StatusUnauthorized)
		invResp := &shared.Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, validErrs
	}
	return nil, nil
}

// VDeleteAccountRequest validate delete account request
func VDeleteAccountRequest(ctx *gin.Context, obj interface{}) (*shared.Response, url.Values) {
	validErrs := url.Values{}
	if invResp, err := shared.ValidateRequest(ctx, obj); err != nil {
		validErrs.Add("error", err.Error())
		return invResp, validErrs
	}
	user := obj.(*Entity)

	if user.Email == "" {
		validErrs.Add("email", "email is required")
	}
	data := make(map[string]interface{})
	if len(validErrs) > 0 {
		data["error"] = validErrs
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
		data["status"] = http.StatusText(http.StatusUnauthorized)
		invResp := &shared.Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, validErrs
	}
	return nil, nil
}

// VListAccountRequest validates list account request
func VListAccountRequest(ctx *gin.Context, obj interface{}) (*shared.Response, url.Values) {
	validErrs := url.Values{}
	if invResp, err := shared.ValidateRequest(ctx, obj); err != nil {
		validErrs.Add("error", err.Error())
		return invResp, validErrs
	}
	user := obj.(*Entity)

	if user.Email == "" {
		validErrs.Add("email", "email is required")
	}
	data := make(map[string]interface{})
	if len(validErrs) > 0 {
		data["error"] = validErrs
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
		data["status"] = http.StatusText(http.StatusUnauthorized)
		invResp := &shared.Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, validErrs
	}
	return nil, nil
}

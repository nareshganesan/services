package account

import (
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"net/http"
)

// SignupUsecase handles the processing for signup request
func SignupUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	query := account.GetAuthQuery()
	if query == nil {
		data["error"] = "Could not create User!"
		data["message"] = "BadReques"
		data["code"] = http.StatusUnauthorized
	} else {
		if account.IsExistingUser(query) {
			data["error"] = "User already exists!"
			data["message"] = "BadRequest"
			data["code"] = http.StatusUnauthorized
		} else {
			id, status := account.Create()
			if !status {
				data["error"] = "Could not create User!"
				data["message"] = "BadRequest"
				data["code"] = http.StatusBadRequest
			} else {
				data["username"] = account.Username
				data["password"] = account.Password
				data["id"] = id
				data["code"] = http.StatusOK
			}
		}
	}
	return shared.GetResponse(ctx, data)
}

// LoginUsecase handles the processing for login request
func LoginUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	query := account.GetAuthQuery()
	if query == nil {
		data["error"] = "Invalid credentials!"
		data["message"] = "StatusUnauthorized"
		data["code"] = http.StatusUnauthorized
	} else {
		if !account.Authenticate(query) {
			data["error"] = "Invalid credentials!"
			data["message"] = "StatusUnauthorized"
			data["code"] = http.StatusBadRequest
		} else {
			token := g.GenerateJWT(account.Username)
			if token == "" {
				data["error"] = "Invalid credentials!"
				data["message"] = "StatusUnauthorized"
				data["code"] = http.StatusBadRequest
			} else {
				data["account"] = account
				data["code"] = http.StatusOK
				data["token"] = token
			}
		}
	}
	return shared.GetResponse(ctx, data)
}

// UpdateAccountUsecase handles the processing for update account request
func UpdateAccountUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	query := account.GetAuthQuery()
	if account.IsExistingUser(query) {
		id, status := account.Update()
		if !status {
			data["error"] = "Could not update User!"
			data["message"] = "BadRequest"
			data["code"] = http.StatusBadRequest
		} else {
			data["username"] = account.Username
			data["password"] = account.Password
			data["id"] = id
			data["code"] = http.StatusOK
		}
	} else {
		data["error"] = "StatusUnauthorized"
		data["message"] = "User does not exist!"
		data["code"] = http.StatusBadRequest
	}
	return shared.GetResponse(ctx, data)
}

// DeleteAccountUsecase handles the processing for delete account request
func DeleteAccountUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	query := account.GetAuthQuery()
	if account.IsExistingUser(query) {
		id, status := account.Delete()
		if !status {
			data["error"] = "Could not delete User!"
			data["message"] = "BadRequest"
			data["code"] = http.StatusBadRequest
		} else {
			data["username"] = account.Username
			data["password"] = account.Password
			data["id"] = id
			data["code"] = http.StatusOK
		}
	} else {
		data["error"] = "StatusUnauthorized"
		data["message"] = "User does not exist!"
		data["code"] = http.StatusBadRequest
	}
	return shared.GetResponse(ctx, data)
}

// ListAccountUsecase handles the processing for list account request
func ListAccountUsecase(ctx *gin.Context, account Entity, page, size int) *shared.Response {
	data := make(map[string]interface{})
	query := account.GetAuthQuery()
	if account.IsExistingUser(query) {
		listquery := GetListAccountQuery()
		accounts := account.List(page, size, listquery)
		if accounts == nil {
			data["error"] = "Could not list Users!"
			data["message"] = "BadRequest"
			data["code"] = http.StatusBadRequest
		} else {
			data["accounts"] = accounts
			data["code"] = http.StatusOK
		}
	} else {
		data["error"] = "StatusUnauthorized"
		data["message"] = "User does not exist!"
		data["code"] = http.StatusBadRequest
	}
	return shared.GetResponse(ctx, data)
}

package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"net/http"
)

// SignupUsecase handles the processing for signup request
func SignupUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	if account.IsExistingUser() {
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
			data["email"] = account.Email
			data["password"] = account.Password
			data["id"] = id
			data["code"] = http.StatusOK
		}
	}
	return shared.GetResponse(ctx, data)
}

// LoginUsecase handles the processing for login request
func LoginUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	if !account.Authenticate() {
		adminEmail := g.Config.Owner.Email
		if account.ID == "" {
			data["message"] = "Invalid credentials!"
		} else {
			if account.IsArchived {
				data["message"] = fmt.Sprintf("Account archived! contact %s to unlock it", adminEmail)
			} else {
				if !account.IsLocked {
					if account.FailedAttempts == 5 {
						data["message"] = fmt.Sprintf("Account locked! use forgot password link, to reset password! (reached maximum failed attempts!)")
					} else {
						data["message"] = fmt.Sprintf("Invalid Credentials!, failed attempt:%d, maximum allowed:%d!", account.FailedAttempts, 5)
					}
				} else {
					data["message"] = fmt.Sprintf("Account locked! use forgot password link, to reset password! (reached maximum failed attempts!)")
				}
			}
		}
		data["error"] = "StatusUnauthorized!"
		data["code"] = http.StatusBadRequest
	} else {
		token := g.GenerateJWT(account.ID)
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
	return shared.GetResponse(ctx, data)
}

// UpdateAccountUsecase handles the processing for update account request
func UpdateAccountUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	account.ID = ctx.GetString("uid")
	account.IsVerified = true
	id, status := account.Update()
	if !status {
		data["error"] = "Could not update User!"
		data["message"] = "BadRequest"
		data["code"] = http.StatusBadRequest
	} else {
		data["email"] = account.Email
		data["password"] = account.Password
		data["id"] = id
		data["code"] = http.StatusOK
	}
	return shared.GetResponse(ctx, data)
}

// DeleteAccountUsecase handles the processing for delete account request
func DeleteAccountUsecase(ctx *gin.Context, account Entity) *shared.Response {
	data := make(map[string]interface{})
	account.ID = ctx.GetString("uid")
	id, status := account.Delete()
	if !status {
		data["error"] = "Could not delete User!"
		data["message"] = "BadRequest"
		data["code"] = http.StatusBadRequest
	} else {
		data["email"] = account.Email
		data["password"] = account.Password
		data["id"] = id
		data["code"] = http.StatusOK
	}
	return shared.GetResponse(ctx, data)
}

// ListAccountUsecase handles the processing for list account request
func ListAccountUsecase(ctx *gin.Context, account Entity, page, size int) *shared.Response {
	data := make(map[string]interface{})
	listquery := GetListAccountQuery()
	account.ID = ctx.GetString("uid")
	accounts := account.List(page, size, listquery)
	if accounts == nil {
		data["error"] = "Could not list Users!"
		data["message"] = "BadRequest"
		data["code"] = http.StatusBadRequest
	} else {
		data["accounts"] = accounts
		data["code"] = http.StatusOK
	}
	return shared.GetResponse(ctx, data)
}

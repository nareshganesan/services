package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
)

// Signup handler for signup request
func Signup(ctx *gin.Context) {
	var account Entity
	if invResp, errs := VRegisterRequest(ctx, &account); len(errs) > 0 {
		invResp.Send()
		return
	}
	resp := SignupUsecase(ctx, account)
	resp.Send()
	return
}

// Login handler for login request
func Login(ctx *gin.Context) {
	Gbl := g.GetGlobals()
	fmt.Println(Gbl)
	var account Entity
	if invResp, errs := VLoginRequest(ctx, &account); len(errs) > 0 {
		invResp.Send()
		return
	}
	resp := LoginUsecase(ctx, account)
	resp.Send()
	return
}

// UpdateAccount handler for update account request
func UpdateAccount(ctx *gin.Context) {
	var account Entity
	if invResp, errs := VUpdateAccountRequest(ctx, &account); len(errs) > 0 {
		invResp.Send()
		return
	}
	resp := UpdateAccountUsecase(ctx, account)
	resp.Send()
	return
}

// DeleteAccount handler for delete account request
func DeleteAccount(ctx *gin.Context) {
	var account Entity
	if invResp, errs := VDeleteAccountRequest(ctx, &account); len(errs) > 0 {
		invResp.Send()
		return
	}
	resp := DeleteAccountUsecase(ctx, account)
	resp.Send()
	return
}

// ListAccount handler for list account request
func ListAccount(ctx *gin.Context) {
	account := Entity{}
	page, size := VListAccountRequest(ctx)
	resp := ListAccountUsecase(ctx, account, page, size)
	resp.Send()
	return
}

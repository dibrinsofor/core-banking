package handlers_test

import (
	"fmt"
	"testing"

	"github.com/dibrinsofor/core-banking/handlers"
	"github.com/jaswdr/faker"
	"gotest.tools/v3/assert"
)

func TestWithdraw200(t *testing.T) {
	f := faker.New()

	req := handlers.MakeTestRequest(t, "/createAccount", map[string]interface{}{
		"name":  f.Person().Name(),
		"email": f.Person().Contact().Email,
	}, "POST")

	verifyResponse := handlers.BootstrapServer(req, routeHandlers)
	verifyResponseBody := handlers.DecodeResponse(t, verifyResponse)
	account_number := verifyResponseBody["data"].(map[string]interface{})["account_number"]
	fmt.Print(account_number)

	depositRequest := handlers.MakeTestRequest(t, "/deposit", map[string]interface{}{
		"account_number": account_number,
		"amount":         20000,
	}, "POST")

	_ = handlers.BootstrapServer(depositRequest, routeHandlers)

	withdrawRequest := handlers.MakeTestRequest(t, "/withdraw", map[string]interface{}{
		"account_number": account_number,
		"amount":         500,
	}, "POST")

	getWithdrawalResponse := handlers.BootstrapServer(withdrawRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getWithdrawalResponse)
	assert.Equal(t, "user withdrawal successful", responseBody["message"])
}

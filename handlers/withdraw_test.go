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
		"name":   f.Person().Name(),
		"email":  f.Person().Contact().Email,
		"amount": 100000,
	}, "POST")

	verifyResponse := handlers.BootstrapServer(req, routeHandlers)
	verifyResponseBody := handlers.DecodeResponse(t, verifyResponse)
	fmt.Print(verifyResponseBody)
	account_number := verifyResponseBody["data"].(map[string]interface{})["account_number"]

	depositRequest := handlers.MakeTestRequest(t, "/withdraw", map[string]interface{}{
		"account_number": account_number,
		"amount":         0,
	}, "POST")

	getDepositResponse := handlers.BootstrapServer(depositRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)
	assert.Equal(t, "user withdrawal successful", responseBody["message"])
}

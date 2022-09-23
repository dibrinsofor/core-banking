package handlers_test

import (
	"testing"

	"github.com/dibrinsofor/core-banking/internal/handlers"
	"github.com/jaswdr/faker"
	"gotest.tools/v3/assert"
)

func TestDeposit200(t *testing.T) {
	f := faker.New()

	req := handlers.MakeTestRequest(t, "/createAccount", map[string]interface{}{
		"name":  f.Person().Name(),
		"email": f.Person().Contact().Email,
	}, "POST")

	verifyResponse := handlers.BootstrapServer(req, routeHandlers)
	verifyResponseBody := handlers.DecodeResponse(t, verifyResponse)
	account_number := verifyResponseBody["data"].(map[string]interface{})["account_number"]

	depositRequest := handlers.MakeTestRequest(t, "/deposit", map[string]interface{}{
		"account_number": account_number,
		"amount":         f.Int16(),
	}, "POST")

	getDepositResponse := handlers.BootstrapServer(depositRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)
	assert.Equal(t, "user deposit successful", responseBody["message"])
}

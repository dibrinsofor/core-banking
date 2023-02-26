package handlers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
		"amount":         fmt.Sprintf("%v", f.Int16()),
	}, "POST")

	getDepositResponse := handlers.BootstrapServer(depositRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)
	assert.Equal(t, "deposit successful", responseBody["message"])
}

func TestDeposit400AmountAsInt(t *testing.T) {
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
	assert.Equal(t, "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD", responseBody["message"])
}

func TestDeposit400InvalidUser(t *testing.T) {
	f := faker.New()

	req := handlers.MakeTestRequest(t, "/createAccount", map[string]interface{}{
		"name":  f.Person().Name(),
		"email": f.Person().Contact().Email,
	}, "POST")

	handlers.BootstrapServer(req, routeHandlers)
	account_number := "not an account number"

	depositRequest := handlers.MakeTestRequest(t, "/deposit", map[string]interface{}{
		"account_number": account_number,
		"amount":         fmt.Sprintf("%v", f.Int16()),
	}, "POST")

	getDepositResponse := handlers.BootstrapServer(depositRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)
	assert.Equal(t, "failed to find user. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD", responseBody["message"])
}

func TestDepositExceedRequestTimeout(t *testing.T) {
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
		"amount":         fmt.Sprintf("%v", f.Int16()),
	}, "POST")

	cancellingCtx, cancel := context.WithCancel(depositRequest.Context())
	time.AfterFunc(5*time.Millisecond, cancel)

	getDepositResponse := handlers.BootstrapServer(depositRequest.WithContext(cancellingCtx), routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)
	assert.Equal(t, "Request timed out.", responseBody["message"])
}

func TestDuplicateDeposit(t *testing.T) {
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
		"amount":         fmt.Sprintf("%v", f.Int16()),
	}, "POST")

	getDepositResponse := handlers.BootstrapServer(depositRequest, routeHandlers)
	responseBody := handlers.DecodeResponse(t, getDepositResponse)

	depositRequest2 := handlers.MakeTestRequest(t, "/deposit", map[string]interface{}{
		"account_number": account_number,
		"amount":         fmt.Sprintf("%v", f.Int16()),
	}, "POST")

	getDepositResponse2 := handlers.BootstrapServer(depositRequest2, routeHandlers)
	responseBody2 := handlers.DecodeResponse(t, getDepositResponse2)

	assert.Equal(t, responseBody2["message"], responseBody["message"])
}

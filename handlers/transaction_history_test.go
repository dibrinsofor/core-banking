package handlers_test

import (
	"fmt"
	"testing"

	"github.com/dibrinsofor/core-banking/handlers"
	"github.com/jaswdr/faker"
	"gotest.tools/v3/assert"
)

func TestTransactionHistory(t *testing.T) {
	f := faker.New()

	// abstract this into function
	user1Req := handlers.MakeTestRequest(t, "/createAccount", map[string]interface{}{
		"name":  f.Person().Name(),
		"email": f.Person().Contact().Email,
	}, "POST")

	verifyUser1Response := handlers.BootstrapServer(user1Req, routeHandlers)
	verifyUser1ResponseBody := handlers.DecodeResponse(t, verifyUser1Response)
	user1account_number := verifyUser1ResponseBody["data"].(map[string]interface{})["account_number"]

	depositRequest := handlers.MakeTestRequest(t, "/deposit", map[string]interface{}{
		"account_number": user1account_number,
		"amount":         20000,
	}, "POST")

	_ = handlers.BootstrapServer(depositRequest, routeHandlers)

	transReq := handlers.MakeTestRequest(t, "/transHistory", map[string]interface{}{
		"account_number": user1account_number,
	}, "GET")

	getTransResponse := handlers.BootstrapServer(transReq, routeHandlers)
	fmt.Printf("%v", getTransResponse.Body)
	responseBody := handlers.DecodeResponse(t, getTransResponse)
	assert.Equal(t, "successfully retrieved 10 most recent transactions", responseBody["message"])
}

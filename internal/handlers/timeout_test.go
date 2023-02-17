package handlers_test

import (
	"testing"

	"github.com/dibrinsofor/core-banking/internal/handlers"
	"gotest.tools/v3/assert"
)

func TestTimeout(t *testing.T) {

	req := handlers.MakeTestRequest(t, "/timeout", map[string]interface{}{}, "GET")

	verifyResponse := handlers.BootstrapServer(req, routeHandlers)
	responseBody := handlers.DecodeResponse(t, verifyResponse)

	assert.Equal(t, "Request timed out.", responseBody["message"])

}

package handlers_test

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/dibrinsofor/core-banking/internal/config"
	"github.com/dibrinsofor/core-banking/internal/handlers"
	"github.com/dibrinsofor/core-banking/internal/postgres"
	"github.com/dibrinsofor/core-banking/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jaswdr/faker"
	"gotest.tools/v3/assert"
)

var routeHandlers *gin.Engine

func TestMain(m *testing.M) {
	err := config.LoadTestConfig("../../.env.test")
	if err != nil {
		panic(err)
	}

	db, err := postgres.Init()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)

	// rdb, err := redistest.Init()
	// if err != nil {
	// 	log.Fatalf("Unable to connect to redis store: %v", err)
	// }

	ht, err := strconv.ParseInt(os.Getenv("HANDLER_TIMEOUT"), 0, 64)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.New(repo, ht)
	srv := config.New(handler)
	routeHandlers = srv.SetupRoutes()
	os.Exit(m.Run())
}

func TestCreateAccountEndpoint(t *testing.T) {
	f := faker.New()

	req := handlers.MakeTestRequest(t, "/createAccount", map[string]interface{}{
		"name":  f.Person().Name(),
		"email": f.Person().Contact().Email,
	}, "POST")

	response := handlers.BootstrapServer(req, routeHandlers)
	responseBody := handlers.DecodeResponse(t, response)

	assert.Equal(t, "user successfully created", responseBody["message"])
}

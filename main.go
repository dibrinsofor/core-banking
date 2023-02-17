package main

import (
	"log"
	"os"
	"strconv"

	"github.com/dibrinsofor/core-banking/internal/config"
	"github.com/dibrinsofor/core-banking/internal/handlers"
	"github.com/dibrinsofor/core-banking/internal/postgres"
	"github.com/dibrinsofor/core-banking/internal/repository"
)

func main() {
	err := config.LoadNormalConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Init()
	if err != nil {
		log.Fatal(err)
	}

	// rdb, err := redistest.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	repo := repository.New(db)

	ht, err := strconv.ParseInt(os.Getenv("HANDLER_TIMEOUT"), 0, 64)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.New(repo, ht)
	server := config.New(handler)

	server.Start()
}

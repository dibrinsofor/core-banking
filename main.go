package main

import (
	"log"

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

	repo := repository.New(db)

	handler := handlers.New(repo)
	server := config.New(handler)

	server.Start()
}

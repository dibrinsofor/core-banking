package handlers

import (
	"time"

	"github.com/dibrinsofor/core-banking/internal/repository"
)

type Handler struct {
	repo            *repository.Repository
	TimeoutDuration time.Duration
}

// add redis client when time comes
func New(repo *repository.Repository, TimeoutDuration int64) *Handler {
	return &Handler{
		repo:            repo,
		TimeoutDuration: time.Duration(time.Duration(TimeoutDuration) * time.Second),
	}
}

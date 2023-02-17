package handlers

import (
	"time"

	"github.com/dibrinsofor/core-banking/internal/repository"
)

type Handler struct {
	repo *repository.Repository
	// rdb             *redis.Client
	TimeoutDuration time.Duration
}

// add redis client when time comes, see tinderclone
func New(repo *repository.Repository, TimeoutDuration int64) *Handler {
	return &Handler{
		repo: repo,
		// rdb:             rdb,
		TimeoutDuration: time.Duration(time.Duration(TimeoutDuration) * time.Second),
	}
}

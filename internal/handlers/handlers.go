package handlers

import (
	"github.com/dibrinsofor/core-banking/internal/repository"
	"github.com/go-redis/redis"
)

type Handler struct {
	repo *repository.Repository
	rdb  *redis.Client
}

// add redis client when time comes, see tinderclone
func New(repo *repository.Repository, rdb *redis.Client) *Handler {
	return &Handler{
		repo: repo,
		rdb:  rdb,
	}
}

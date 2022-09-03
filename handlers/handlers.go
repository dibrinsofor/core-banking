package handlers

import (
	"github.com/dibrinsofor/core-banking/repository"
)

type Handler struct {
	repo *repository.Repository
}

// add redis client when time comes, see tinderclone
func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

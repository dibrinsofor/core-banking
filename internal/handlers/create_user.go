package handlers

import (
	"log"
	"net/http"

	"github.com/dibrinsofor/core-banking/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateUserPayload struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser CreateUserPayload
	var v models.AccountInfo

	err := c.BindJSON(&newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	v.Email = newUser.Email
	v.Name = newUser.Name
	v.AccountBalance = "$0"

	if err = h.repo.UserRepo.CreateUser(&v); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user successfully created",
		"data":    v,
	})
}

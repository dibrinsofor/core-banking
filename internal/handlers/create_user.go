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
	var v models.Users

	err := c.BindJSON(&newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	balance, err := h.repo.UserRepo.GetUserBalance("")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find user",
		})
		return
	}

	v.Email = newUser.Email
	v.Name = newUser.Name
	v.Balance = 0

	s := SanitizeAmount(balance)
	userBalance := s.(string)

	if err = h.repo.UserRepo.CreateUser(&v); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user successfully created",
		"data": gin.H{
			"account_number": v.AccountNumber,
			"name":           v.Name,
			"email":          v.Email,
			"balance":        userBalance,
			"created_at":     v.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

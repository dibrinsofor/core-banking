package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepositRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	Amount        int    `json:"deposit" binding:"required"`
}

func (h *Handler) Deposit(c *gin.Context) {
	var updateUser DepositRequest

	if err := c.BindJSON(&updateUser); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	existingUserData, err := h.repo.UserRepo.GetUserByID(updateUser.AccountNumber)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find user",
		})
		return
	}

	existingUserData.Balance = existingUserData.Balance + updateUser.Amount

	err = h.repo.UserRepo.UpdateUserByID(updateUser.AccountNumber, existingUserData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deposit successful",
		"data":    existingUserData,
	})
}

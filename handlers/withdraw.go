package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WithdrawRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}

func (h *Handler) Withdraw(c *gin.Context) {
	var updateUser WithdrawRequest

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

	if existingUserData.Balance < updateUser.Amount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "balance too low to process payment",
		})
	}

	existingUserData.Balance = existingUserData.Balance - updateUser.Amount

	err = h.repo.UserRepo.UpdateUserByID(updateUser.AccountNumber, existingUserData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, something went wrong. failed to perform transaction.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user withdrawal successful",
		"data":    existingUserData,
	})
}

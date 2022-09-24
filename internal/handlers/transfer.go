package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	Recipient     string `json:"recipient" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
}

func (h *Handler) Transfer(c *gin.Context) {
	var transferReq TransferRequest

	if err := c.BindJSON(&transferReq); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	existingUserData, err := h.repo.UserRepo.GetUserByID(transferReq.AccountNumber)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find user. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	recipientData, err := h.repo.UserRepo.GetUserByID(transferReq.Recipient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find user. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	i := SanitizeAmount(transferReq.Amount)
	amount := i.(float64)

	if existingUserData.Balance < amount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "balance too low to process payment. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
	}

	newBalance, err := h.repo.UserRepo.Transfer(existingUserData, recipientData, amount)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, failed to complete transfer. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transfer successful",
		"data": gin.H{
			"account_number": existingUserData.AccountNumber,
			"name":           existingUserData.Name,
			"balance":        newBalance,
			"recipient_name": recipientData.Name,
			"updated_at":     time.Now().Format("2017-09-07 17:06:06"),
		},
	})
}

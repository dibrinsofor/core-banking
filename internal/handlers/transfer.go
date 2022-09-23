package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	AccountNumber string  `json:"account_number" binding:"required"`
	Recipient     string  `json:"recipient" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
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
			"message": "failed to find user",
		})
		return
	}

	recipientData, err := h.repo.UserRepo.GetUserByID(transferReq.Recipient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find user",
		})
		return
	}

	if existingUserData.Balance < transferReq.Amount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "balance too low to process payment",
		})
	}

	// bundle db changes into a pipeline if something fails nothing should be processed
	existingUserData.Balance = existingUserData.Balance - transferReq.Amount
	recipientData.Balance = recipientData.Balance + transferReq.Amount

	err = h.repo.UserRepo.UpdateUsersByID(existingUserData.AccountNumber, recipientData.AccountNumber, existingUserData, recipientData, "TRANSFER")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, something went wrong. failed to perform transaction.",
		})
		return
	}

	// maybe just send acc number and updated balance
	c.JSON(http.StatusOK, gin.H{
		"message": "transfer successful",
		"data":    existingUserData,
	})
}

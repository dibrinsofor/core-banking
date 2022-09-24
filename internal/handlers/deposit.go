package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DepositRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
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
			"message": "failed to find user. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	i := SanitizeAmount(updateUser.Amount)
	updateUserAmount := i.(float64)

	balance, err := h.repo.UserRepo.Deposit(existingUserData, updateUserAmount)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, somethingwent wrong. Failed to make deposit. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	s := SanitizeAmount(&balance)
	userBalance := s.(string)

	c.JSON(http.StatusOK, gin.H{
		"message": "deposit successful",
		"data": gin.H{
			"account_number": updateUser.AccountNumber,
			"name":           existingUserData.Name,
			"balance":        userBalance,
			"updated_at":     time.Now().Format("2017-09-07 17:06:06"),
		},
	})
}

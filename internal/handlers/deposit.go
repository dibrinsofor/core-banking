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
	s := SanitizeAmount(balance)
	userBalance := s.(string)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, something went wrong. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	uData := gin.H{
		"message": "deposit successful",
		"data": gin.H{
			"account_number": updateUser.AccountNumber,
			"name":           existingUserData.Name,
			"balance":        userBalance,
			"updated_at":     time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// key := c.Request.Header.Get(redistest.DefaultKeyName)

	// userInstance := &redistest.Idempotency{
	// 	KeyName:         key,
	// 	CreatedAt:       time.Now(),
	// 	CacheExpiration: time.Now().Add(time.Hour * 2),
	// 	SavedResponse:   uData,
	// }

	// if key != "" {
	// 	redistest.AddIdempotencyKey(userInstance)
	// }

	c.JSON(http.StatusOK, uData)
}

package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestHistory struct {
	AccountNumber string `json:"account_number" binding:"required"`
}

type QueryParam struct {
	Action string
	Date   time.Time
}

func (h *Handler) TransactionHistory(c *gin.Context) {
	var reqHistory RequestHistory

	if err := c.BindJSON(&reqHistory); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	action := c.DefaultQuery("action", "")
	date := c.Query("date")
	// dateTime, err := time.Parse("2006-01-02", date)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "failed to parse date string",
	// 	})
	// 	return
	// }

	if action == "" && date == "" {
		trans, err := h.repo.UserRepo.GetAllTransactions(reqHistory.AccountNumber)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "uhoh, something went wrong. failed to fetch transactions.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    trans,
			"message": "successfully retrieved 10 most recent transactions",
		})
		return
	}

	trans, err := h.repo.UserRepo.QueryTransactions(reqHistory.AccountNumber, action, date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "uhoh, something went wrong. failed to fetch transactions.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully retrieved transactions",
		"data":    trans,
	})
}

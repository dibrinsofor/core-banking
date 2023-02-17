package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Timeout(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success. Request completed before timeout.",
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "up and running")
}

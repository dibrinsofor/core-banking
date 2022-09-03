package handlers

import (
	"log"
	"net/http"

	"github.com/dibrinsofor/core-banking/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser models.User

	err := c.BindJSON(&newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse user request. check documentation: https://github.com/dibrinsofor/core-banking/blob/master/Readme.MD",
		})
		return
	}

	if err = h.repo.UserRepo.CreateUser(&newUser); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user successfully created",
		"data":    newUser,
	})
}

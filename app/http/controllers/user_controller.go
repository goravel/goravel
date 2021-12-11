package controllers

import "github.com/gin-gonic/gin"

type UserController struct {
}

func (user UserController) Show(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "show",
	})
}

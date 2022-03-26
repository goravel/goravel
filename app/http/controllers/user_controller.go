package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
)

type UserController struct {
}

func (r UserController) Show(ctx *gin.Context) {
	facades.Response.Success(ctx, gin.H{
		"Hello": "Goravel",
	})
}

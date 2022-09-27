package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/goravel/framework/facades"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx *gin.Context) {
	facades.Response.Success(ctx, gin.H{
		"Hello": "Goravel",
	})
}

package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Aakash-Sleur/go-micro-auth/models"
	"github.com/Aakash-Sleur/go-micro-auth/service"
	"github.com/Aakash-Sleur/go-micro-auth/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.AuthService
}

func NewUserController(servive *service.AuthService) *UserController {
	return &UserController{service: servive}
}

func (c *UserController) SignUp(ctx *gin.Context) {
	var req struct {
		Name     string
		Email    string
		Password string
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := c.service.SignUp(req.Email, req.Name, req.Password)
	if err != nil {
		if err.Error() == "email already in use" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User created Successfully", gin.H{"token": token, "user": user})
}

func (c *UserController) Signin(ctx *gin.Context) {
	var req struct {
		Email    string
		Password string
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := c.service.Signin(req.Email, req.Password)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User created Successfully", gin.H{"token": token, "user": user})
}

func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	userInterface, exists := ctx.Get("currentUser")

	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not in context").Error())
		return
	}

	fmt.Println(userInterface, "userinterface")

	user := userInterface.(models.User)

	utils.SuccessResponse(ctx, http.StatusOK, "User created Successfully", gin.H{"id": user.ID, "email": user.Email})
}

package controller

import (
	commonModels "commonpkg/models"
	"net/http"
	"organization-service/service"

	"github.com/gin-gonic/gin"
)

var userControllerObj *UserController

type UserController struct {
	svc *service.UserService
}

func InitUserController() (*UserController, *commonModels.ErrorDetail) {
	if userControllerObj == nil {
		service, err := service.InitUserService()
		if err != nil {
			return nil, err
		}
		userControllerObj = &UserController{
			svc: service,
		}
	}
	return userControllerObj, nil
}

func (ctrl *UserController) Login(context *gin.Context) {
	var loginRequest commonModels.LoginRequest
	if err := context.ShouldBindJSON(&loginRequest); err == nil {
		data := ctrl.svc.Login(loginRequest)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

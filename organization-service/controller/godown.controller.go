package controller

import (
	commonModels "commonpkg/models"
	"net/http"
	"organization-service/service"

	"github.com/gin-gonic/gin"
)

var godownControllerObj *GodownController

type GodownController struct {
	svc *service.GodownService
}

func InitGodownController() (*GodownController, *commonModels.ErrorDetail) {
	if godownControllerObj == nil {
		service, err := service.InitGodownService()
		if err != nil {
			return nil, err
		}
		godownControllerObj = &GodownController{
			svc: service,
		}
	}
	return godownControllerObj, nil
}

func (ctrl *GodownController) GetAll(context *gin.Context) {
	data := ctrl.svc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *GodownController) Add(context *gin.Context) {
	var addData commonModels.GodownAddRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.svc.Add(addData.Name)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

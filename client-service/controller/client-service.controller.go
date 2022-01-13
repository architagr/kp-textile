package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"client-service/service"

	"github.com/gin-gonic/gin"
)

var clientServiceCtr *ClientServiceController

type ClientServiceController struct {
	clientServiceSvc *service.ClientServiceService
}

func InitClientServiceController() (*ClientServiceController, *commonModels.ErrorDetail) {
	if clientServiceCtr == nil {
		svc, err := service.InitClientServiceService()
		if err != nil {
			return nil, err
		}
		clientServiceCtr = &ClientServiceController{
			clientServiceSvc: svc,
		}
	}
	return clientServiceCtr, nil
}

func (ctrl *ClientServiceController) Add(context *gin.Context) {
	var addData commonModels.AddClientRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.clientServiceSvc.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"transportor-service/service"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	Code string 
}

type AddMultipleRequest struct {
	Codes []string 
}
type GetRequest struct {
	Id string 
}

var transportorServiceCtr *TransportorServiceController

type TransportorServiceController struct {
	transportorServiceSvc *service.TransportorServiceService
}

func InitTransportorServiceController() (*TransportorServiceController, *commonModels.ErrorDetail) {
	if transportorServiceCtr == nil {
		svc, err := service.InitTransportorServiceService()
		if err != nil {
			return nil, err
		}
		transportorServiceCtr = &TransportorServiceController{
			transportorServiceSvc: svc,
		}
	}
	return transportorServiceCtr, nil
}
func (ctrl *TransportorServiceController) GetAll(context *gin.Context) {
	data := ctrl.transportorServiceSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *TransportorServiceController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.transportorServiceSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *TransportorServiceController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.transportorServiceSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *TransportorServiceController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.transportorServiceSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

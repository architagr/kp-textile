package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"item-service/service"

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

var itemServiceCtr *ItemServiceController

type ItemServiceController struct {
	itemServiceSvc *service.ItemServiceService
}

func InitItemServiceController() (*ItemServiceController, *commonModels.ErrorDetail) {
	if itemServiceCtr == nil {
		svc, err := service.InitItemServiceService()
		if err != nil {
			return nil, err
		}
		itemServiceCtr = &ItemServiceController{
			itemServiceSvc: svc,
		}
	}
	return itemServiceCtr, nil
}
func (ctrl *ItemServiceController) GetAll(context *gin.Context) {
	data := ctrl.itemServiceSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *ItemServiceController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.itemServiceSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *ItemServiceController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.itemServiceSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *ItemServiceController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.itemServiceSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

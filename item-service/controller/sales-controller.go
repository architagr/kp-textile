package controller

import (
	commonModels "commonpkg/models"
	"item-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var salesControllerObj *SalesController

type SalesController struct {
	salesService *service.SalesService
}

func InitSalesController() (*SalesController, *commonModels.ErrorDetail) {
	if salesControllerObj == nil {
		svc, err := service.InitSalesService()
		if err != nil {
			return nil, err
		}

		salesControllerObj = &SalesController{
			salesService: svc,
		}
	}

	return salesControllerObj, nil
}

func (ctrl *SalesController) GetAllSalesOrders(context *gin.Context) {
	var getAllRequest commonModels.InventoryListRequest
	if err := context.ShouldBindJSON(&getAllRequest); err == nil {
		data := ctrl.salesService.GetAll(getAllRequest)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *SalesController) AddSalesBillDetails(context *gin.Context) {
	var addData commonModels.AddSalesDataRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.salesService.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

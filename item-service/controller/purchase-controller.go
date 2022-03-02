package controller

import (
	commonModels "commonpkg/models"
	"item-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var purchaseControllerObj *PurchaseController

type PurchaseController struct {
	purchaseService service.IPurchaseService
}

func InitPurchaseController() (*PurchaseController, *commonModels.ErrorDetail) {
	if purchaseControllerObj == nil {
		svc, err := service.InitPurchaseService()
		if err != nil {
			return nil, err
		}

		purchaseControllerObj = &PurchaseController{
			purchaseService: svc,
		}
	}

	return purchaseControllerObj, nil
}

func (ctrl *PurchaseController) GetAllPurchaseOrders(context *gin.Context) {
	var getAllRequest commonModels.InventoryListRequest
	if err := context.ShouldBindJSON(&getAllRequest); err == nil {
		data := ctrl.purchaseService.GetAll(getAllRequest)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *PurchaseController) AddPurchaseBillDetails(context *gin.Context) {
	var addData commonModels.AddPurchaseDataRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.purchaseService.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var purchaseControllerObj *PurchaseController

type PurchaseController struct {
	purchaseService *service.PurchaseService
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
		getAllRequest.BranchId = getBranchIdFromContext(context)
		data := ctrl.purchaseService.GetAllPurchaseOrders(getAllRequest)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *PurchaseController) GetPurchaseBillDetails(context *gin.Context) {
	var request commonModels.InventoryFilterDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext((context))

		data := ctrl.purchaseService.GetPurchaseBillDetails(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *PurchaseController) AddPurchaseBillDetails(context *gin.Context) {
	var addData commonModels.InventoryDto
	if err := context.ShouldBindJSON(&addData); err == nil {
		addData.BranchId = getBranchIdFromContext(context)
		data := ctrl.purchaseService.AddPurchaseBillDetails(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *PurchaseController) UpdatePurchaseBillDetails(context *gin.Context) {
	var updateData commonModels.InventoryDto
	var filterData commonModels.InventoryFilterDto
	if err := context.ShouldBindJSON(&updateData); err == nil {
		if err1 := context.ShouldBindUri(&filterData); err1 == nil {

			updateData.BillNo = filterData.PurchaseBillNumber
			updateData.BranchId = getBranchIdFromContext(context)
			data := ctrl.purchaseService.UpdatePurchaseBillDetails(updateData)
			context.JSON(data.StatusCode, data)
		} else {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err1,
			})
		}
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *PurchaseController) DeletePurchaseBillDetails(context *gin.Context) {
	var filterData commonModels.InventoryFilterDto
	if err := context.ShouldBindUri(&filterData); err == nil {
		filterData.BranchId = getBranchIdFromContext(context)
		data := ctrl.purchaseService.DeletePurchaseBillDetails(filterData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func getBranchIdFromContext(context *gin.Context) string {
	return fmt.Sprint(context.Keys[commonModels.ContextKey_BranchId])
}

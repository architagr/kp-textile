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
		getAllRequest.BranchId = getBranchIdFromContext(context)
		data := ctrl.salesService.GetAllSalesOrders(getAllRequest)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *SalesController) GetSalesBillDetails(context *gin.Context) {
	var request commonModels.InventoryFilterDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext((context))

		data := ctrl.salesService.GetSalesBillDetails(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *SalesController) AddSalesBillDetails(context *gin.Context) {
	var addData commonModels.InventoryDto
	if err := context.ShouldBindJSON(&addData); err == nil {
		addData.BranchId = getBranchIdFromContext(context)
		data := ctrl.salesService.AddSalesBillDetails(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *SalesController) UpdateSalesBillDetails(context *gin.Context) {
	var updateData commonModels.InventoryDto
	var filterData commonModels.InventoryFilterDto
	if err := context.ShouldBindJSON(&updateData); err == nil {
		if err1 := context.ShouldBindUri(&filterData); err1 == nil {
			updateData.BillNo = filterData.SalesBillNumber
			updateData.BranchId = getBranchIdFromContext(context)
			data := ctrl.salesService.UpdateSalesBillDetails(updateData)
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

func (ctrl *SalesController) DeleteSalesBillDetails(context *gin.Context) {
	var filterData commonModels.InventoryFilterDto
	if err := context.ShouldBindUri(&filterData); err == nil {
		filterData.BranchId = getBranchIdFromContext(context)
		data := ctrl.salesService.DeleteSalesBillDetails(filterData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

package controller

import (
	commonModels "commonpkg/models"

	"github.com/gin-gonic/gin"
)

var salesControllerObj *SalesController

type SalesController struct {
}

func InitSalesController() (*SalesController, *commonModels.ErrorDetail) {
	if salesControllerObj == nil {
		salesControllerObj = &SalesController{}
	}

	return salesControllerObj, nil
}

func (ctrl *SalesController) GetAllSalesOrders(context *gin.Context) {

}

func (ctrl *SalesController) GetSalesBillDetails(context *gin.Context) {

}

func (ctrl *SalesController) AddSalesBillDetails(context *gin.Context) {

}

func (ctrl *SalesController) UpdateSalesBillDetails(context *gin.Context) {

}

func (ctrl *SalesController) DeleteSalesBillDetails(context *gin.Context) {

}

package controller

import (
	commonModels "commonpkg/models"
	"item-service/service"
)

var baleControllerObj *BaleController

type BaleController struct {
	baleService *service.BaleService
}

func InitBaleController() (*BaleController, *commonModels.ErrorDetail) {
	if baleControllerObj == nil {
		svc, err := service.InitBaleService()
		if err != nil {
			return nil, err
		}
		baleControllerObj = &BaleController{
			baleService: svc,
		}
	}

	return baleControllerObj, nil
}

// func (ctrl *BaleController) GetBaleInfo(context *gin.Context) {
// 	var filterData commonModels.BaleInfoReuest
// 	if err := context.ShouldBindUri(&filterData); err == nil {
// 		data := ctrl.baleService.GetBaleInfo(filterData)
// 		context.JSON(data.StatusCode, data)
// 	} else {
// 		context.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 	}
// }

// func (ctrl *BaleController) GetBalesByQuantity(context *gin.Context) {
// 	var filterData commonModels.BaleInfoReuest
// 	if err := context.ShouldBindUri(&filterData); err == nil {
// 		data := ctrl.baleService.GetBaleInfoByQuality(filterData)
// 		context.JSON(data.StatusCode, data)
// 	} else {
// 		context.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 	}
// }

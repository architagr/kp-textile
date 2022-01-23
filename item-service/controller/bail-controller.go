package controller

import (
	commonModels "commonpkg/models"
	"item-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var bailControllerObj *BailController

type BailController struct {
	bailService *service.BailService
}

func InitBailController() (*BailController, *commonModels.ErrorDetail) {
	if bailControllerObj == nil {
		svc, err := service.InitBailService()
		if err != nil {
			return nil, err
		}
		bailControllerObj = &BailController{
			bailService: svc,
		}
	}

	return bailControllerObj, nil
}

func (ctrl *BailController) GetBailInfo(context *gin.Context) {
	var filterData commonModels.BailInfoReuest
	if err := context.ShouldBindUri(&filterData); err == nil {
		filterData.BranchId = getBranchIdFromContext(context)
		data := ctrl.bailService.GetBailInfo(filterData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *BailController) GetBailsByQuantity(context *gin.Context) {
	var filterData commonModels.BailInfoReuest
	if err := context.ShouldBindUri(&filterData); err == nil {
		filterData.BranchId = getBranchIdFromContext(context)
		data := ctrl.bailService.GetBailInfoByQuality(filterData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

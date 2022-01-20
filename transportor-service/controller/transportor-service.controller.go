package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"transportor-service/service"

	"github.com/gin-gonic/gin"
)

var transporterCtr *TransporterController

type TransporterController struct {
	transporterServiceSvc *service.TransporterService
}

func InitTransporterController() (*TransporterController, *commonModels.ErrorDetail) {
	if transporterCtr == nil {
		svc, err := service.InitTransporterService()
		if err != nil {
			return nil, err
		}
		transporterCtr = &TransporterController{
			transporterServiceSvc: svc,
		}
	}
	return transporterCtr, nil
}

func (ctrl *TransporterController) Add(context *gin.Context) {
	var addData commonModels.AddTransporterRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		addData.BranchId = getBranchIdFromContext(context)
		data := ctrl.transporterServiceSvc.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *TransporterController) Get(context *gin.Context) {
	var request commonModels.GetTransporterRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext(context)

		data := ctrl.transporterServiceSvc.GetTransporter(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *TransporterController) Delete(context *gin.Context) {
	var request commonModels.GetTransporterRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext(context)

		data := ctrl.transporterServiceSvc.DeleteTransporter(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}
func (ctrl *TransporterController) Put(context *gin.Context) {
	var request commonModels.AddTransporterRequest

	if err := context.ShouldBindJSON(&request); err == nil {
		if err1 := context.ShouldBindUri(&request); err1 == nil {
			request.BranchId = getBranchIdFromContext(context)

			data := ctrl.transporterServiceSvc.Put(request)
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

func (ctrl *TransporterController) GetAll(context *gin.Context) {
	var getAllRequest commonModels.TransporterListRequest
	if err := context.ShouldBindJSON(&getAllRequest); err == nil {
		if err1 := context.ShouldBindQuery(&getAllRequest); err1 == nil {
			getAllRequest.BranchId = getBranchIdFromContext((context))
			data := ctrl.transporterServiceSvc.GetAll(getAllRequest)
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

func getBranchIdFromContext(context *gin.Context) string {
	return fmt.Sprint(context.Keys[commonModels.ContextKey_BranchId])
}

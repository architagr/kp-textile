package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"client-service/service"

	"github.com/gin-gonic/gin"
)

var clientServiceCtr *ClientServiceController

type ClientServiceController struct {
	clientServiceSvc *service.ClientServiceService
}

func InitClientServiceController() (*ClientServiceController, *commonModels.ErrorDetail) {
	if clientServiceCtr == nil {
		svc, err := service.InitClientServiceService()
		if err != nil {
			return nil, err
		}
		clientServiceCtr = &ClientServiceController{
			clientServiceSvc: svc,
		}
	}
	return clientServiceCtr, nil
}

func (ctrl *ClientServiceController) Add(context *gin.Context) {
	var addData commonModels.AddClientRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		addData.BranchId = getBranchIdFromContext(context)
		data := ctrl.clientServiceSvc.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *ClientServiceController) Get(context *gin.Context) {
	var request commonModels.GetClientRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext((context))

		data := ctrl.clientServiceSvc.GetClient(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *ClientServiceController) Delete(context *gin.Context) {
	var request commonModels.GetClientRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext(context)

		data := ctrl.clientServiceSvc.DeleteClient(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}
func (ctrl *ClientServiceController) Put(context *gin.Context) {
	var request commonModels.AddClientRequest

	if err := context.ShouldBindJSON(&request); err == nil {
		if err1 := context.ShouldBindUri(&request); err1 == nil {
			request.BranchId = getBranchIdFromContext(context)

			data := ctrl.clientServiceSvc.Put(request)
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

func (ctrl *ClientServiceController) GetAll(context *gin.Context) {
	var getAllRequest commonModels.ClientListRequest
	if err := context.ShouldBindJSON(&getAllRequest); err == nil {
		if err1 := context.ShouldBindQuery(&getAllRequest); err1 == nil {
			getAllRequest.BranchId = getBranchIdFromContext((context))
			data := ctrl.clientServiceSvc.GetAll(getAllRequest)
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

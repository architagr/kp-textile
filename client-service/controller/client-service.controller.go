package controller

import (
	commonModels "commonpkg/models"
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
		//TODO: add branch id from middleware
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
		//TODO: add branch id from middleware
		request.BranchId = "branchId"

		data := ctrl.clientServiceSvc.GetClient(request)
		context.JSON(data.StatusCode, data)
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
			//TODO: add branch id from middleware
			getAllRequest.BranchId = "branchId"
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

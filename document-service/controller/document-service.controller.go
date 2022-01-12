package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"document-service/service"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	Code string 
}

type AddMultipleRequest struct {
	Codes []string 
}
type GetRequest struct {
	Id string 
}

var documentServiceCtr *DocumentServiceController

type DocumentServiceController struct {
	documentServiceSvc *service.DocumentServiceService
}

func InitDocumentServiceController() (*DocumentServiceController, *commonModels.ErrorDetail) {
	if documentServiceCtr == nil {
		svc, err := service.InitDocumentServiceService()
		if err != nil {
			return nil, err
		}
		documentServiceCtr = &DocumentServiceController{
			documentServiceSvc: svc,
		}
	}
	return documentServiceCtr, nil
}
func (ctrl *DocumentServiceController) GetAll(context *gin.Context) {
	data := ctrl.documentServiceSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *DocumentServiceController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.documentServiceSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *DocumentServiceController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.documentServiceSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *DocumentServiceController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.documentServiceSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

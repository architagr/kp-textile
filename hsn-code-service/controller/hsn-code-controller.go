package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"hsn-code-service/service"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	Code string `json:"code" binding:"required"`
}

type AddMultipleRequest struct {
	Codes []string `json:"codes"`
}
type GetRequest struct {
	Id string `uri:"id"`
}

var hnsCodeCtr *HnsCodeController

type HnsCodeController struct {
	hnsCodeSvc *service.HnsCodeService
}

func InitHnsCodeController() (*HnsCodeController, *commonModels.ErrorDetail) {
	if hnsCodeCtr == nil {
		svc, err := service.InitHnsCodeService()
		if err != nil {
			return nil, err
		}
		hnsCodeCtr = &HnsCodeController{
			hnsCodeSvc: svc,
		}
	}
	return hnsCodeCtr, nil
}
func (ctrl *HnsCodeController) GetAll(context *gin.Context) {
	data := ctrl.hnsCodeSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *HnsCodeController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.hnsCodeSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *HnsCodeController) Add(context *gin.Context) {
	var addData AddRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.hnsCodeSvc.Add(addData.Code)
		fmt.Println("data received ")
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *HnsCodeController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.hnsCodeSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"quality-service/service"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	Code string `json:"code" binding:"required"`
}

type AddMultipleRequest struct {
	Codes []string `json:"codes" binding:"required"`
}
type GetRequest struct {
	Id string `uri:"id"`
}

var qualityCtr *QualityController

type QualityController struct {
	qualitySvc *service.QualityService
}

func InitQualityController() (*QualityController, *commonModels.ErrorDetail) {
	if qualityCtr == nil {
		svc, err := service.InitQualityService()
		if err != nil {
			return nil, err
		}
		qualityCtr = &QualityController{
			qualitySvc: svc,
		}
	}
	return qualityCtr, nil
}
func (ctrl *QualityController) GetAll(context *gin.Context) {
	data := ctrl.qualitySvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *QualityController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.qualitySvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *QualityController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.qualitySvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *QualityController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.qualitySvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

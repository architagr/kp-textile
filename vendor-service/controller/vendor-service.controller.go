package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"vendor-service/service"

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

var vendorServiceCtr *VendorServiceController

type VendorServiceController struct {
	vendorServiceSvc *service.VendorServiceService
}

func InitVendorServiceController() (*VendorServiceController, *commonModels.ErrorDetail) {
	if vendorServiceCtr == nil {
		svc, err := service.InitVendorServiceService()
		if err != nil {
			return nil, err
		}
		vendorServiceCtr = &VendorServiceController{
			vendorServiceSvc: svc,
		}
	}
	return vendorServiceCtr, nil
}
func (ctrl *VendorServiceController) GetAll(context *gin.Context) {
	data := ctrl.vendorServiceSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *VendorServiceController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.vendorServiceSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *VendorServiceController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.vendorServiceSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *VendorServiceController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.vendorServiceSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

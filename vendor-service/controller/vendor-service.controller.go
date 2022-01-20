package controller

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"vendor-service/service"

	"github.com/gin-gonic/gin"
)

var vendorCtr *VendorController

type VendorController struct {
	vendorServiceSvc *service.VendorService
}

func InitVendorController() (*VendorController, *commonModels.ErrorDetail) {
	if vendorCtr == nil {
		svc, err := service.InitVendorService()
		if err != nil {
			return nil, err
		}
		vendorCtr = &VendorController{
			vendorServiceSvc: svc,
		}
	}
	return vendorCtr, nil
}

func (ctrl *VendorController) Add(context *gin.Context) {
	var addData commonModels.AddVendorRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		addData.BranchId = getBranchIdFromContext(context)
		data := ctrl.vendorServiceSvc.Add(addData)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *VendorController) Get(context *gin.Context) {
	var request commonModels.GetVendorRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext(context)

		data := ctrl.vendorServiceSvc.GetVendor(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *VendorController) Delete(context *gin.Context) {
	var request commonModels.GetVendorRequestDto

	if err := context.ShouldBindUri(&request); err == nil {
		request.BranchId = getBranchIdFromContext(context)

		data := ctrl.vendorServiceSvc.DeleteVendor(request)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}
func (ctrl *VendorController) Put(context *gin.Context) {
	var request commonModels.AddVendorRequest

	if err := context.ShouldBindJSON(&request); err == nil {
		if err1 := context.ShouldBindUri(&request); err1 == nil {
			request.BranchId = getBranchIdFromContext(context)

			data := ctrl.vendorServiceSvc.Put(request)
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

func (ctrl *VendorController) GetAll(context *gin.Context) {
	var getAllRequest commonModels.VendorListRequest
	if err := context.ShouldBindJSON(&getAllRequest); err == nil {
		if err1 := context.ShouldBindQuery(&getAllRequest); err1 == nil {
			getAllRequest.BranchId = getBranchIdFromContext((context))
			data := ctrl.vendorServiceSvc.GetAll(getAllRequest)
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

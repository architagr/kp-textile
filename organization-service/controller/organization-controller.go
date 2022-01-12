package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"organization-service/service"

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

var organizationCtr *OrganizationController

type OrganizationController struct {
	organizationSvc *service.OrganizationService
}

func InitOrganizationController() (*OrganizationController, *commonModels.ErrorDetail) {
	if organizationCtr == nil {
		svc, err := service.InitHnsCodeService()
		if err != nil {
			return nil, err
		}
		organizationCtr = &OrganizationController{
			organizationSvc: svc,
		}
	}
	return organizationCtr, nil
}
func (ctrl *OrganizationController) GetAll(context *gin.Context) {
	data := ctrl.organizationSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *OrganizationController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.organizationSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *OrganizationController) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.organizationSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *OrganizationController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.organizationSvc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

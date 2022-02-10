package controller

import (
	commonModels "commonpkg/models"
	"net/http"

	"quality-service/service"

	"github.com/gin-gonic/gin"
)

type AddMultipleRequest struct {
	Names []string `json:"names" binding:"required"`
}

var productCtr *ProductController

type ProductController struct {
	productSvc *service.ProductService
}

func InitProductController() (*ProductController, *commonModels.ErrorDetail) {
	if productCtr == nil {
		svc, err := service.InitProductService()
		if err != nil {
			return nil, err
		}
		productCtr = &ProductController{
			productSvc: svc,
		}
	}
	return productCtr, nil
}
func (ctrl *ProductController) GetAll(context *gin.Context) {
	data := ctrl.productSvc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *ProductController) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.productSvc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

}

func (ctrl *ProductController) Add(context *gin.Context) {
	var addData AddRequest

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.productSvc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

func (ctrl *ProductController) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.productSvc.AddMultiple(addData.Names)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}

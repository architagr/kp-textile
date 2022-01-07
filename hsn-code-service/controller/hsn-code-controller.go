package controller

import (
	"hsn-code-service/model"
	"hsn-code-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	code string `form:"code"`
}

type AddMultipleRequest struct {
	codes []string `form:"codes"`
}

var hnsCodeCtr *HnsCodeController

// TODO: analysis response from service and in case of error send error response
type HnsCodeController struct {
	hnsCodeSvc *service.HnsCodeService
}

func InitHnsCodeController(svc *service.HnsCodeService) *HnsCodeController {
	if hnsCodeCtr == nil {
		return &HnsCodeController{
			hnsCodeSvc: svc,
		}
	}

	return hnsCodeCtr
}
func (ctrl *HnsCodeController) GetAll(context *gin.Context) []model.HnsCodeDto {
	return ctrl.hnsCodeSvc.GetAll()
}

func (ctrl *HnsCodeController) Get(context *gin.Context) model.HnsCodeDto {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

	}
	return ctrl.hnsCodeSvc.Get(id)
}

func (ctrl *HnsCodeController) Add(context *gin.Context) model.HnsCodeDto {
	var addData AddRequest
	if context.ShouldBind(&addData) == nil {
		return ctrl.hnsCodeSvc.Add(addData.code)
	} else {
		// FIXME: error
		return model.HnsCodeDto{}
	}
}

func (ctrl *HnsCodeController) AddMultiple(context *gin.Context) []model.HnsCodeDto {
	var addData AddMultipleRequest
	if context.ShouldBind(&addData) == nil {
		return ctrl.hnsCodeSvc.AddMultiple(addData.codes)
	} else {
		// FIXME: error
		return []model.HnsCodeDto{}

	}
}

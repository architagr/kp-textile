package controller

import (
	"hsn-code-service/model"
	"hsn-code-service/service"

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

func InitHnsCodeController() *HnsCodeController {
	if hnsCodeCtr == nil {
		return &HnsCodeController{
			hnsCodeSvc: service.InitHnsCodeService(),
		}
	}

	return hnsCodeCtr
}
func (ctrl *HnsCodeController) GetAll(context *gin.Context) {
	data := ctrl.hnsCodeSvc.GetAll()
	context.JSON(200, data)
}

func (ctrl *HnsCodeController) Get(context *gin.Context) {
	id := context.Param("id")
	data := ctrl.hnsCodeSvc.Get(id)
	context.JSON(200, data)

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

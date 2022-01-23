package controller

import (
	commonModels "commonpkg/models"

	"github.com/gin-gonic/gin"
)

var bailControllerObj *BailController

type BailController struct {
}

func InitBailController() (*BailController, *commonModels.ErrorDetail) {
	if bailControllerObj == nil {
		bailControllerObj = &BailController{}
	}

	return bailControllerObj, nil
}

func (ctrl *BailController) GetBailInfo(context *gin.Context) {

}

func (ctrl *BailController) CheckBailNumber(context *gin.Context) {

}

func (ctrl *BailController) GetBailsByQuantity(context *gin.Context) {

}

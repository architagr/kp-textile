package controller

import (
	commonModels "commonpkg/models"
	"time"

	"document-service/service"

	"github.com/gin-gonic/gin"
)

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

func (ctrl *DocumentServiceController) GetChallan(context *gin.Context) {
	var data commonModels.InventoryDto = commonModels.InventoryDto{
		BranchId:         "branchId",
		InventorySortKey: "Inventory|Deleted|sales01|Inventory|Sales|1643142960",
		PurchaseDate:     time.Date(2022, 12, 25, 0, 0, 0, 0, nil),
		SalesDate:        time.Date(2022, 12, 25, 0, 0, 0, 0, nil),
		BillNo:           "sales01",
		LrNo:             "123",
		ChallanNo:        "123",
		HsnCode:          "sdsdf",
		BailDetails: []commonModels.BailDetailsDto{
			{
				BailNo:         "bail01",
				BilledQuantity: 1000,
				Quality:        "4729adb5-7432-11ec-a804-0800275114e0",
			},
		},
		TransporterId: "daac2ed7-7a3f-11ec-bda7-0800275114e0",
	}
	var htmlBody = `<html>
	<head>
	<style>
	@media print
	{
	  table { page-break-after:auto }
	  tr    { page-break-inside:avoid; page-break-after:auto }
	  td    { page-break-inside:avoid; page-break-after:auto }
	  thead { display:table-header-group }
	  tfoot { display:table-footer-group }
	}
	</style>
	</head>
	
	<body>`

	htmlBody = htmlBody + `<table>`
}

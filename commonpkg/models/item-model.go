package models

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type PurchaseDto struct {
	GodownId       string    `json:"godownId,omitempty"`
	SortKey        string    `json:"sortKey,omitempty"`        /// ProductId|QualityId|PurchaseId
	PurchaseId     string    `json:"purchaseId,omitempty"`     // GSI -1  PK  (all attr)
	PurchaseBillNo string    `json:"purchaseBillNo,omitempty"` // GSI - 2 PK (keys only)
	Date           time.Time `json:"date,omitempty"`
	VendorId       string    `json:"vendorId,omitempty"`
	ProductId      string    `json:"productId,omitempty"`
	QualityId      string    `json:"qualityId,omitempty"`
	Status         string    `json:"status,omitempty"` // Stock | Sold
}

type SalesDto struct {
	GodownId      string    `json:"godownId,omitempty"`
	SortKey       string    `json:"sortKey,omitempty"`     /// ProductId|QualityId
	SalesId       string    `json:"salesId,omitempty"`     // GSI -1  PK  (all attr)
	SalesBillNo   string    `json:"salesBillNo,omitempty"` // GSI PK (keys only)
	ClientId      string    `json:"clientId,omitempty"`
	TransporterId string    `json:"transporterId,omitempty"`
	LrNo          string    `json:"lrNo,omitempty"`
	ChallanNo     string    `json:"challanNo,omitempty"` // GSI PK (keys only)
	Date          time.Time `json:"date,omitempty"`
	ProductId     string    `json:"productId,omitempty"`
	QualityId     string    `json:"qualityId,omitempty"`
}

type BalePurchaseDetails struct {
	PurchaseId string
}
type BaleSalesDetails struct {
	SalesId string
}
type BaleTransferDetails struct {
	FromGodownId string    `json:"fromGodownId,omitempty"`
	ToGowodnId   string    `json:"toGowodnId,omitempty"`
	Date         time.Time `json:"date,omitempty"`
}

type BaleDetailsDto struct {
	GodownId         string                `json:"godownId,omitempty"`
	SortKey          string                `json:"sortKey,omitempty"` //// <InStock or OutOfStock>|ProductId|QualityId|BaleNo
	BaleNo           string                `json:"baleNo,omitempty"`  //GSI PK (all attr)
	ProductId        string                `json:"productId,omitempty"`
	QualityId        string                `json:"qualityId,omitempty"`
	BilledQuantity   int32                 `json:"billedQuantity,omitempty"`
	ReceivedQuantity int32                 `json:"receivedQuantity,omitempty"`
	Rate             int32                 `json:"rate,omitempty"`
	PurchaseDetails  BalePurchaseDetails   `json:"purchaseDetails,omitempty"`
	SalesDetails     BaleSalesDetails      `json:"salesDetails,omitempty"`
	TransferDetails  []BaleTransferDetails `json:"transferDetails,omitempty"`
}

type InventoryFilterDto struct {
	GodownId           string    `json:"godownId,omitempty"`
	PurchaseBillNumber string    `json:"purchaseBillNumber,omitempty" uri:"purchaseBillNumber"`
	SalesBillNumber    string    `json:"salesBillNumber,omitempty" uri:"salesBillNumber"`
	StartDate          time.Time `json:"startDate,omitempty"`
	EndDate            time.Time `json:"endDate,omitempty"`
	QualityId          string    `json:"qualityId,omitempty"`
	ProductId          string    `json:"productId,omitempty"`
}

type InventoryListRequest struct {
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty" form:"pageSize"`
	PurchaseId       string
	PurchaseBillNo   string
	SalesId          string
	InventoryFilterDto
}

type PurchaseListResponse struct {
	CommonListResponse
	Data []PurchaseDto `json:"data,omitempty"`
}
type PurchaseResponse struct {
	CommonResponse
	Data PurchaseDto `json:"data,omitempty"`
}

type AddPurchaseDataRequest struct {
	PurchaseDetails PurchaseDto      `json:"purchaseDetails,omitempty"`
	BaleDetails     []BaleDetailsDto `json:"baleDetails,omitempty"`
}
type AddPurchaseDataResponse struct {
	CommonResponse
	PurchaseDetails PurchaseDto      `json:"purchaseDetails,omitempty"`
	BaleDetails     []BaleDetailsDto `json:"baleDetails,omitempty"`
}

// type InventoryResponse struct {
// 	CommonResponse
// 	Data InventoryDto `json:"data,omitempty"`
// }

// type BaleInfoReuest struct {
// 	GodownId string
// 	BaleNo   string `uri:"baleNo,omitempty"`
// 	Quality  string `uri:"quality,omitempty"`
// }

// type BaleInfoResponse struct {
// 	CommonResponse
// 	Purchase []BaleDetailsDto `json:"purchase,omitempty"`
// 	Sales    []BaleDetailsDto `json:"sales,omitempty"`
// }

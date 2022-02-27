package models

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type PurchaseDto struct {
	GodownId       string    `json:"godownId,omitempty"`
	SortKey        string    `json:"sortKey,omitempty"`        /// ProductId|QualityId
	PurchaseId     string    `json:"purchaseId,omitempty"`     // GSI -1  PK  (all attr)
	PurchaseBillNo string    `json:"purchaseBillNo,omitempty"` // GSI - 2 PK (keys only)
	Date           time.Time `json:"date,omitempty"`
	VendorId       string    `json:"vendorId,omitempty"`
	ProductId      string    `json:"productId,omitempty"`
	QualityId      string    `json:"qualityId,omitempty"`
	Status         string    `json:"status,omitempty"`
}

type SalesDto struct {
	GodownId      string    `json:"godownId,omitempty"`
	SortKey       string    `json:"sortKey,omitempty"` /// ProductId|QualityId
	SalesId       string    `json:"salesId,omitempty"`
	SalesBillNo   string    `json:"salesBillNo,omitempty"` // GSI PK
	ClientId      string    `json:"clientId,omitempty"`
	TransporterId string    `json:"transporterId,omitempty"`
	LrNo          string    `json:"lrNo,omitempty"`
	ChallanNo     string    `json:"challanNo,omitempty"` // GSI PK
	Date          time.Time `json:"date,omitempty"`
	ProductId     string    `json:"productId,omitempty"`
	QualityId     string    `json:"qualityId,omitempty"`
	Status        string    `json:"status,omitempty"`
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
	SortKey          string                `json:"sortKey,omitempty"` //// <Stock or OutOfStock>|ProductId|QualityId|BaleNo
	BaleNo           string                `json:"baleNo,omitempty"`  //GSI PK
	ProductId        string                `json:"productId,omitempty"`
	QualityId        string                `json:"qualityId,omitempty"`
	BilledQuantity   int32                 `json:"billedQuantity,omitempty"`
	ReceivedQuantity int32                 `json:"receivedQuantity,omitempty"`
	Rate             int32                 `json:"rate,omitempty"`
	PurchaseDetails  BalePurchaseDetails   `json:"purchaseDetails,omitempty"`
	SalesDetails     BaleSalesDetails      `json:"salesDetails,omitempty"`
	TransferDetails  []BaleTransferDetails `json:"transferDetails,omitempty"`
}

// type BaleInfoDto struct {
// 	GodownId         string `json:"godownId,omitempty"`
// 	BaleInfoSortKey  string `json:"baleInfoSortKey,omitempty"` /// Info | baleNo | Quality
// 	BaleNo           string `json:"baleNo,omitempty"`
// 	ReceivedQuantity int32  `json:"receivedQuantity,omitempty"`
// 	BilledQuantity   int32  `json:"billedQuantity,omitempty"`
// 	IsLongation      bool   `json:"isLongation,omitempty"`
// 	Quality          string `json:"quality,omitempty"`
// }

// type InventoryDto struct {
// 	GodownId         string           `json:"godownId,omitempty"`
// 	InventorySortKey string           `json:"inventorySortKey,omitempty"` /// Inventory | <Purchase or Sales>| Bill No
// 	BillNo           string           `json:"billNo,omitempty"`
// 	BaleDetails      []BaleDetailsDto `json:"baleDetails,omitempty"`
// 	PurchaseDate     time.Time        `json:"purchaseDate,omitempty" time_format:"2006-01-02"`
// 	SalesDate        time.Time        `json:"salesDate,omitempty" time_format:"unix"`
// 	TransporterId    string           `json:"transporterId,omitempty"`
// 	LrNo             string           `json:"lrNo,omitempty"`
// 	ChallanNo        string           `json:"challanNo,omitempty"`
// 	HsnCode          string           `json:"hsnCode,omitempty"`
// }

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
	SalesId          string
	InventoryFilterDto
}

// type InventoryListResponse struct {
// 	CommonListResponse
// 	Data []InventoryDto `json:"data,omitempty"`
// }

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

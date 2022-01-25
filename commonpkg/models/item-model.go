package models

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type BailDetailsDto struct {
	BranchId               string    `json:"branchId,omitempty"`
	SortKey                string    `json:"sortKey,omitempty"` //// Bail | <Purchase or Sales or OutOfStock>|Quality|BailNo|<salesBill or purchseBill number>
	BailNo                 string    `json:"bailNo,omitempty"`
	Quality                string    `json:"quality,omitempty"`
	IsSales                bool      `json:"isSales,omitempty"`
	BillNo                 string    `json:"billNo,omitempty"`
	Rate                   int32     `json:"rate,omitempty"`
	PurchaseDate           time.Time `json:"purchaseDate,omitempty"`
	SalesDate              time.Time `json:"salesDate,omitempty"`
	ClientId               string    `json:"clientId,omitempty"`
	VendorId               string    `json:"vendorId,omitempty"`
	TransferedToBranchId   string    `json:"transferedToBranchId,omitempty"`
	TransferedFromBranchId string    `json:"transferedFromBranchId,omitempty"`
	ReceivedQuantity       int32     `json:"receivedQuantity,omitempty"`
	BilledQuantity         int32     `json:"billedQuantity,omitempty"`
	PendingQuantity        int32     `json:"pendingQuantity,omitempty"`
}

type BailInfoDto struct {
	BranchId         string `json:"branchId,omitempty"`
	BailInfoSortKey  string `json:"bailInfoSortKey,omitempty"` /// Info | bailNo | Quality
	BailNo           string `json:"bailNo,omitempty"`
	ReceivedQuantity int32  `json:"receivedQuantity,omitempty"`
	BilledQuantity   int32  `json:"billedQuantity,omitempty"`
	IsLongation      bool   `json:"isLongation,omitempty"`
	Quality          string `json:"quality,omitempty"`
}

type InventoryDto struct {
	BranchId         string           `json:"branchId,omitempty"`
	InventorySortKey string           `json:"inventorySortKey,omitempty"` /// Inventory | <Purchase or Sales>| Bill No
	BillNo           string           `json:"billNo,omitempty"`
	BailDetails      []BailDetailsDto `json:"bailDetails,omitempty"`
	PurchaseDate     time.Time        `json:"purchaseDate,omitempty" time_format:"2006-01-02"`
	SalesDate        time.Time        `json:"salesDate,omitempty" time_format:"unix"`
	TransporterId    string           `json:"transporterId,omitempty"`
	LrNo             string           `json:"lrNo,omitempty"`
	ChallanNo        string           `json:"challanNo,omitempty"`
	HsnCode          string           `json:"hsnCode,omitempty"`
}

type InventoryFilterDto struct {
	BranchId           string
	PurchaseBillNumber string    `json:"purchaseBillNumber,omitempty" uri:"purchaseBillNumber"`
	SalesBillNumber    string    `json:"salesBillNumber,omitempty" uri:"salesBillNumber"`
	StartDate          time.Time `json:"startDate,omitempty"`
	EndDate            time.Time `json:"endDate,omitempty"`
	Quality            string    `json:"quality,omitempty"`
}

type InventoryListRequest struct {
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty" form:"pageSize"`
	InventoryFilterDto
}
type InventoryListResponse struct {
	CommonListResponse
	Data []InventoryDto `json:"data,omitempty"`
}

type InventoryResponse struct {
	CommonResponse
	Data InventoryDto `json:"data,omitempty"`
}

type BailInfoReuest struct {
	BranchId string
	BailNo   string `uri:"bailNo,omitempty"`
	Quality  string `uri:"quality,omitempty"`
}

type BailInfoResponse struct {
	CommonResponse
	Purchase []BailDetailsDto `json:"purchase,omitempty"`
	Sales    []BailDetailsDto `json:"sales,omitempty"`
}

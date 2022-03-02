package common

import (
	"os"
)

type Env struct {
	BaleTableName           string
	PurchaseTableName       string
	SalesTableName          string
	PurchaseIdIndexName     string
	SalesIdIndexName        string
	PurchaseBillNoIndexName string
	SalesBillNoIndexName    string
	ChallanNoIndexName      string
	BaleNoIndexName         string
	PortNumber              string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			BaleTableName:           os.Getenv("BaleTable"),
			PurchaseTableName:       os.Getenv("PurchaseTable"),
			SalesTableName:          os.Getenv("SalesTable"),
			PurchaseIdIndexName:     os.Getenv("PurchaseIdIndex"),
			SalesIdIndexName:        os.Getenv("SalesIdIndex"),
			PurchaseBillNoIndexName: os.Getenv("PurchaseBillNoIndex"),
			SalesBillNoIndexName:    os.Getenv("SalesBillNoIndex"),
			ChallanNoIndexName:      os.Getenv("ChallanNoIndex"),
			BaleNoIndexName:         os.Getenv("BaleNoIndex"),
			PortNumber:              "0",
		}
	} else {
		EnvValues = Env{
			BaleTableName:           "bale-table",
			PurchaseTableName:       "purchase-table",
			SalesTableName:          "sales-table",
			SalesIdIndexName:        "salesid-index",
			PurchaseIdIndexName:     "purchaseid-index",
			PurchaseBillNoIndexName: "purchase-billno-index",
			SalesBillNoIndexName:    "sales-billno-index",
			ChallanNoIndexName:      "challanno-index",
			BaleNoIndexName:         "baleno-index",
			PortNumber:              ":8084",
		}
	}
}

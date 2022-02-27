package common

import "os"

type Env struct {
	BaleTableName           string
	PurchaseTableName       string
	SalesTableName          string
	PurchaseIdIndexName     string
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
			PurchaseBillNoIndexName: os.Getenv("PurchaseBillNoIndex"),
			SalesBillNoIndexName:    os.Getenv("SalesBillNoIndex"),
			ChallanNoIndexName:      os.Getenv("ChallanNoIndex"),
			BaleNoIndexName:         os.Getenv("BaleNoIndex"),
			PortNumber:              "0",
		}
	} else {
		EnvValues = Env{
			BaleTableName:           os.Getenv("bale-table"),
			PurchaseTableName:       os.Getenv("purchase-table"),
			SalesTableName:          os.Getenv("sales-table"),
			PurchaseIdIndexName:     os.Getenv("purchaseid-index"),
			PurchaseBillNoIndexName: os.Getenv("purchase-billno-index"),
			SalesBillNoIndexName:    os.Getenv("sales-billno-index"),
			ChallanNoIndexName:      os.Getenv("challanno-index"),
			BaleNoIndexName:         os.Getenv("baleno-index"),
			PortNumber:              ":8084",
		}
	}
}

package common

import "os"

type Env struct {
	ItemTableName      string
	InventoryTableName string
	BailInfoTableName  string
	PortNumber         string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			ItemTableName:      os.Getenv("ItemTable"),
			BailInfoTableName:  os.Getenv("BailInfoTable"),
			InventoryTableName: os.Getenv("InventoryTable"),
			PortNumber:         "0",
		}
	} else {
		EnvValues = Env{
			ItemTableName:      "item-table",
			BailInfoTableName:  "bail-info-table",
			InventoryTableName: "inventory-table",
			PortNumber:         ":8084",
		}
	}
}

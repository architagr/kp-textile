package common

import "os"

type Env struct {
	QualityTableName string
	ProductTableName string
	PortNumber       string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			QualityTableName: os.Getenv("qualityTable"),
			ProductTableName: os.Getenv("productTable"),
			PortNumber:       "0",
		}
	} else {
		EnvValues = Env{
			QualityTableName: "quality-table",
			ProductTableName: "product-table",
			PortNumber:       ":8080",
		}
	}
}

package common

import "os"

type Env struct {
	VendorTableName string
	PortNumber      string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			VendorTableName: os.Getenv("VendorTable"),
			PortNumber:      "0",
		}
	} else {
		EnvValues = Env{
			VendorTableName: "vendor-table",
			PortNumber:      ":8082",
		}
	}
}

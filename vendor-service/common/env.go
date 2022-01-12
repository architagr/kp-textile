package common

import "os"

type Env struct {
	TableName  string
	PortNumber string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			TableName:  os.Getenv("VendorService"),
			PortNumber: "0",
		}
	} else {
		EnvValues = Env{
			TableName:  "vendor-service",
			PortNumber: ":8080",
		}
	}
}

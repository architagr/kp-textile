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
			TableName:  os.Getenv("ItemService"),
			PortNumber: "0",
		}
	} else {
		EnvValues = Env{
			TableName:  "item-service",
			PortNumber: ":8080",
		}
	}
}

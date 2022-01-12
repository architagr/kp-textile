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
			TableName:  os.Getenv("DocumentService"),
			PortNumber: "0",
		}
	} else {
		EnvValues = Env{
			TableName:  "document-service",
			PortNumber: ":8080",
		}
	}
}

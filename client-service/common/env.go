package common

import "os"

type Env struct {
	ClientTableName string
	PortNumber      string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			ClientTableName: os.Getenv("ClientTable"),
			PortNumber:      "0",
		}
	} else {
		EnvValues = Env{
			ClientTableName: "client-table",
			PortNumber:      ":8080",
		}
	}
}

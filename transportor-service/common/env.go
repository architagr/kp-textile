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
			TableName:  os.Getenv("TransportorService"),
			PortNumber: "0",
		}
	} else {
		EnvValues = Env{
			TableName:  "transportor-service",
			PortNumber: ":8080",
		}
	}
}

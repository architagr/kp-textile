package common

import "os"

type Env struct {
	ClientTableName  string
	ContactTableName string
	PortNumber       string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			ClientTableName:  os.Getenv("ClientTable"),
			ContactTableName: os.Getenv("ClientContactTable"),
			PortNumber:       "0",
		}
	} else {
		EnvValues = Env{
			ClientTableName:  "client-table",
			ContactTableName: "client-contact-table",
			PortNumber:       ":8080",
		}
	}
}

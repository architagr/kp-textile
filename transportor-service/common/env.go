package common

import "os"

type Env struct {
	TransporterTableName string
	PortNumber           string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			TransporterTableName: os.Getenv("TransporterTable"),
			PortNumber:           "0",
		}
	} else {
		EnvValues = Env{
			TransporterTableName: "transporter-table",
			PortNumber:           ":8083",
		}
	}
}

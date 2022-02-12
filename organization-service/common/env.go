package common

import "os"

type Env struct {
	GodownTableName string
	UserTableName   string
	PortNumber      string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			GodownTableName: os.Getenv("godownTable"),
			UserTableName:   os.Getenv("userTable"),
			PortNumber:      "0",
		}
	} else {
		EnvValues = Env{
			GodownTableName: "godown-table",
			UserTableName:   "user-table",
			PortNumber:      ":8087",
		}
	}
}

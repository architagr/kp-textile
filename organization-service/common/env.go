package common

import "os"

type Env struct {
	BranchTableName string
	UserTableName   string
	PortNumber      string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			BranchTableName: os.Getenv("branchTable"),
			UserTableName:   os.Getenv("userTable"),
			PortNumber:      "0",
		}
	} else {
		EnvValues = Env{
			BranchTableName: "branch",
			UserTableName:   "User",
			PortNumber:      ":8080",
		}
	}
}

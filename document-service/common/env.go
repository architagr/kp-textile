package common

type Env struct {
	PortNumber string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == "" {
		EnvValues = Env{
			PortNumber: "0",
		}
	} else {
		EnvValues = Env{
			PortNumber: ":8085",
		}
	}
}

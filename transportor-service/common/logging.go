package common

import (
	"commonpkg/customLog"
	"log"
	"os"
)

var logginObj customLog.ILogger

func InitLogger() {
	var err error
	logginObj, err = customLog.Init(1, ":transporter-service: ", os.Stdout)

	if err != nil {
		log.Fatal("logger not initilized")
	}
}

func WriteLog(logLevel int, logMessage string) {
	logginObj.Write(logLevel, logMessage)
}

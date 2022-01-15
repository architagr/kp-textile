package router

import (
	"client-service/common"
	"client-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitClientServiceController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}

	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	engine.GET("/:clientId", func(c *gin.Context) {
		controller.Get(c)
	})

	engine.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

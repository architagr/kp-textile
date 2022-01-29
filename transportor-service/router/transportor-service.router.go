package router

import (
	"commonpkg/middlewares"
	"transportor-service/common"
	"transportor-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitTransporterController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())

	transporterGroup := engine.Group("transporter")
	transporterGroup.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	transporterGroup.GET("/:transporterId", func(c *gin.Context) {
		controller.Get(c)
	})
	transporterGroup.DELETE("/:transporterId", func(c *gin.Context) {
		controller.Delete(c)
	})
	transporterGroup.PUT("/:transporterId", func(c *gin.Context) {
		controller.Put(c)
	})

	transporterGroup.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

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

	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	engine.GET("/:transporterId", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.DELETE("/:transporterId", func(c *gin.Context) {
		controller.Delete(c)
	})
	engine.PUT("/:transporterId", func(c *gin.Context) {
		controller.Put(c)
	})

	engine.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

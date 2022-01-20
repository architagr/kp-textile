package router

import (
	"commonpkg/middlewares"
	"vendor-service/common"
	"vendor-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitVendorController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())

	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	engine.GET("/:vendorId", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.DELETE("/:vendorId", func(c *gin.Context) {
		controller.Delete(c)
	})
	engine.PUT("/:vendorId", func(c *gin.Context) {
		controller.Put(c)
	})

	engine.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

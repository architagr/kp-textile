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

	group := engine.Group("/vendor")
	group.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	group.GET("/:vendorId", func(c *gin.Context) {
		controller.Get(c)
	})
	group.DELETE("/:vendorId", func(c *gin.Context) {
		controller.Delete(c)
	})
	group.PUT("/:vendorId", func(c *gin.Context) {
		controller.Put(c)
	})

	group.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

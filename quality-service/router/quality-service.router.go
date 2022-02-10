package router

import (
	"commonpkg/middlewares"
	"quality-service/common"
	"quality-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	qualityController, err := controller.InitQualityController()
	productController, err := controller.InitProductController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())
	qualityGroup := engine.Group("quality")
	qualityGroup.GET("/", func(c *gin.Context) {
		qualityController.GetAll(c)
	})
	qualityGroup.GET("/:id", func(c *gin.Context) {
		qualityController.Get(c)
	})
	qualityGroup.POST("/", func(c *gin.Context) {
		qualityController.Add(c)
	})
	qualityGroup.POST("/addmultiple", func(c *gin.Context) {
		qualityController.AddMultiple(c)
	})

	productGroup := engine.Group("product")
	productGroup.GET("/", func(c *gin.Context) {
		productController.GetAll(c)
	})
	productGroup.GET("/:id", func(c *gin.Context) {
		productController.Get(c)
	})
	productGroup.POST("/", func(c *gin.Context) {
		productController.Add(c)
	})
	productGroup.POST("/addmultiple", func(c *gin.Context) {
		productController.AddMultiple(c)
	})
}

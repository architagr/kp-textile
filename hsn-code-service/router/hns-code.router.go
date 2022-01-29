package router

import (
	"commonpkg/middlewares"
	"hsn-code-service/common"
	"hsn-code-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitHnsCodeController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())
	hsnCodeGroup := engine.Group("hsncode")
	hsnCodeGroup.GET("/", func(c *gin.Context) {
		controller.GetAll(c)
	})
	hsnCodeGroup.GET("/:id", func(c *gin.Context) {
		controller.Get(c)
	})
	hsnCodeGroup.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})
	hsnCodeGroup.POST("/addmultiple", func(c *gin.Context) {
		controller.AddMultiple(c)
	})
}

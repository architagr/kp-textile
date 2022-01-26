package router

import (
	"commonpkg/middlewares"
	"document-service/common"
	"document-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitDocumentServiceController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())

	engine.POST("/challan", func(c *gin.Context) {
		controller.GetChallan(c)
	})
	// engine.GET("/:id", func(c *gin.Context) {
	// 	controller.Get(c)
	// })
	// engine.POST("/", func(c *gin.Context) {
	// 	controller.Add(c)
	// })
	// engine.POST("/addmultiple", func(c *gin.Context) {
	// 	controller.AddMultiple(c)
	// })
}

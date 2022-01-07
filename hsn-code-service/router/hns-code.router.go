package router

import (
	"hsn-code-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine, controller *controller.HnsCodeController) {
	engine.GET("/", func(c *gin.Context) {
		controller.GetAll(c)
	})
	engine.GET("/:id", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})
	engine.POST("/multiple", func(c *gin.Context) {
		controller.AddMultiple(c)
	})
}

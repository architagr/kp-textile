package router

import (
	"client-service/common"
	"client-service/controller"
	"commonpkg/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitClientServiceController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())
	clientGroup := engine.Group("/client")
	clientGroup.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	clientGroup.GET("/:clientId", func(c *gin.Context) {
		controller.Get(c)
	})
	clientGroup.DELETE("/:clientId", func(c *gin.Context) {
		controller.Delete(c)
	})
	clientGroup.PUT("/:clientId", func(c *gin.Context) {
		controller.Put(c)
	})
	clientGroup.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

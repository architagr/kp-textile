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
	engine.Use(middlewares.CORSMiddleware())
	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})

	engine.GET("/:clientId", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.DELETE("/:clientId", func(c *gin.Context) {
		controller.Delete(c)
	})
	engine.PUT("/:clientId", func(c *gin.Context) {
		controller.Put(c)
	})

	engine.POST("/getall", func(c *gin.Context) {
		controller.GetAll(c)
	})
}

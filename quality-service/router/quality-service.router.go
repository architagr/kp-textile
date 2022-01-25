package router

import (
	"commonpkg/middlewares"
	"quality-service/common"
	"quality-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.InitQualityController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())
	engine.GET("/", func(c *gin.Context) {
		controller.GetAll(c)
	})
	engine.GET("/:id", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})
	engine.POST("/addmultiple", func(c *gin.Context) {
		controller.AddMultiple(c)
	})
}

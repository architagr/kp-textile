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
	qualityGroup := engine.Group("quality")
	qualityGroup.GET("/", func(c *gin.Context) {
		controller.GetAll(c)
	})
	qualityGroup.GET("/:id", func(c *gin.Context) {
		controller.Get(c)
	})
	qualityGroup.POST("/", func(c *gin.Context) {
		controller.Add(c)
	})
	qualityGroup.POST("/addmultiple", func(c *gin.Context) {
		controller.AddMultiple(c)
	})
}

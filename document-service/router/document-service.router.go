package router

import (
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


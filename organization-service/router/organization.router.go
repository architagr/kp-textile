package router

import (
	"commonpkg/middlewares"
	"organization-service/common"
	"organization-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {
	godownController, err := controller.InitGodownController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	userController, err := controller.InitUserController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.Use(middlewares.CORSMiddleware())
	godownGroup := engine.Group("godown")
	godownGroup.Use(middlewares.ValidateTokenMiddleware())
	godownGroup.GET("/", func(c *gin.Context) {
		godownController.GetAll(c)
	})
	godownGroup.POST("/", func(c *gin.Context) {
		godownController.Add(c)
	})

	userGroup := engine.Group("user")
	userGroup.POST("/login", func(c *gin.Context) {
		userController.Login(c)
	})
}

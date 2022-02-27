package router

import (
	"commonpkg/middlewares"
	"item-service/common"
	"item-service/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(engine *gin.Engine) {

	engine.Use(middlewares.CORSMiddleware(), middlewares.ValidateTokenMiddleware())
	purchaseRoutes(engine)
	salesRoutes(engine)
	baleInfo(engine)
}

func baleInfo(engine *gin.Engine) {
	baleControllerObj, err := controller.InitBaleController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}

	baleInfoApiGroup := engine.Group("/baleInfo")

	baleInfoApiGroup.GET("/quality/:quality", func(c *gin.Context) {
		baleControllerObj.GetBalesByQuantity(c)
	})

	baleInfoApiGroup.GET("/:baleNo", func(c *gin.Context) {
		baleControllerObj.GetBaleInfo(c)
	})
}

func salesRoutes(engine *gin.Engine) {
	salesControllerObj, err := controller.InitSalesController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	salesApiGroup := engine.Group("/sales")

	salesApiGroup.POST("/getall", func(c *gin.Context) {
		salesControllerObj.GetAllSalesOrders(c)
	})

	salesApiGroup.GET("/:salesBillNumber", func(c *gin.Context) {
		salesControllerObj.GetSalesBillDetails(c)
	})

	salesApiGroup.POST("/", func(c *gin.Context) {
		salesControllerObj.AddSalesBillDetails(c)
	})
}
func purchaseRoutes(engine *gin.Engine) {

	purchaseControllerObj, err := controller.InitPurchaseController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}

	purchaseApiGroup := engine.Group("/purchase")

	purchaseApiGroup.POST("/getall", func(c *gin.Context) {
		purchaseControllerObj.GetAllPurchaseOrders(c)
	})

	purchaseApiGroup.GET("/:purchaseBillNumber", func(c *gin.Context) {
		purchaseControllerObj.GetPurchaseBillDetails(c)
	})

	purchaseApiGroup.POST("/", func(c *gin.Context) {
		purchaseControllerObj.AddPurchaseBillDetails(c)
	})
}

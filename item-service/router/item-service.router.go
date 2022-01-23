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
	bailInfo(engine)
}

func bailInfo(engine *gin.Engine) {
	bailControllerObj, err := controller.InitBailController()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}

	bailInfoApiGroup := engine.Group("/bailInfo")

	bailInfoApiGroup.GET("/quality/:quality", func(c *gin.Context) {
		bailControllerObj.GetBailsByQuantity(c)
	})

	bailInfoApiGroup.GET("/:bailNo", func(c *gin.Context) {
		bailControllerObj.GetBailInfo(c)
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

	salesApiGroup.PUT("/:salesBillNumber", func(c *gin.Context) {
		salesControllerObj.UpdateSalesBillDetails(c)
	})

	salesApiGroup.DELETE("/:salesBillNumber", func(c *gin.Context) {
		salesControllerObj.DeleteSalesBillDetails(c)
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

	purchaseApiGroup.PUT("/:purchaseBillNumber", func(c *gin.Context) {
		purchaseControllerObj.UpdatePurchaseBillDetails(c)
	})

	purchaseApiGroup.DELETE("/:purchaseBillNumber", func(c *gin.Context) {
		purchaseControllerObj.DeletePurchaseBillDetails(c)

	})
}

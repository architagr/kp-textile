package main

import (
	"context"
	"fmt"
	"hsn-code-service/common"
	"os"

	"hsn-code-service/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

// TODO: init all needed services
// TODO: create a custom error in a common repo
// TODO: move all models to common project

var ginLambda *ginadapter.GinLambda
var ginEngine *gin.Engine
var isLocal string

func init() {
	isLocal = os.Getenv("isLocal")
	fmt.Printf("isLocal %s -\n", isLocal)
	common.InitLogger()

	common.WriteLog(1, "Service start")
	ginEngine = gin.Default()
	router.InitRoutes(ginEngine)
	if isLocal == "" {
		ginLambda = ginadapter.New(ginEngine)
	}
}

// Handler is the function that executes for every Request passed into the Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if isLocal == "" {
		lambda.Start(Handler)
	} else {
		ginEngine.Run()
	}
}

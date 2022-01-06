package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	logger1 := log.Default()
	//gin.SetMode(gin.ReleaseMode)
	logger1.Printf("start %s", "service")
	fmt.Println("fmt print")
	os.Setenv("test-key", "Archit")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"env":     os.Getenv("test-key"),
			"env1":    os.Getenv("test"),
		})
	})
	ginLambda = ginadapter.New(r)
	//r.Run()
}

// Handler is the function that executes for every Request passed into the Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

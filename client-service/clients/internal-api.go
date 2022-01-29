package main

import (
	"client-service/common"
	"client-service/router"
	commonModels "commonpkg/models"
	"commonpkg/token"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda
var ginEngine *gin.Engine
var isLocal string

func init() {
	isLocal = os.Getenv("isLocal")
	common.InitLogger()
	common.InitEnv(isLocal)

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
	timeNow := time.Now()
	expireTime := timeNow.AddDate(0, 2, 0)
	claims := commonModels.JwtClaims{
		BranchId: "branchId",
		Username: "Username",
		Roles:    []int{1},
		StandardClaims: jwt.StandardClaims{
			IssuedAt: timeNow.Unix(),
		},
	}

	toekenStr, _ := token.GenrateToken(&claims, expireTime)

	fmt.Println(toekenStr)
	if isLocal == "" {
		lambda.Start(Handler)
	} else {
		ginEngine.Run(common.EnvValues.PortNumber)
	}
}

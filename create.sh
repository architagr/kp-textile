fileName="client-service"
name="ClientService"
privatename="clientService"

cd "$fileName"

mkdir clients common controller persistance router service
touch ./clients/external-api.go  ./clients/internal-api.go
touch ./common/env.go ./common/logging.go
touch "./controller/$fileName.controller.go"
touch "./persistance/$fileName.persistance.go"
touch "./router/$fileName.router.go"
touch "./service/$fileName.service.go"
go mod init "$fileName"

echo "replace commonpkg => ./../commonpkg" >> go.mod
echo "package main

import (
	\"context\"
	\"$fileName/common\"
	\"os\"

	\"$fileName/router\"

	\"github.com/aws/aws-lambda-go/events\"
	\"github.com/aws/aws-lambda-go/lambda\"
	ginadapter \"github.com/awslabs/aws-lambda-go-api-proxy/gin\"
	\"github.com/gin-gonic/gin\"
)

var ginLambda *ginadapter.GinLambda
var ginEngine *gin.Engine
var isLocal string

func init() {
	isLocal = os.Getenv(\"isLocal\")
	common.InitLogger()
	common.InitEnv(isLocal)

	common.WriteLog(1, \"Service start\")
	ginEngine = gin.Default()
	router.InitRoutes(ginEngine)
	if isLocal == \"\" {
		ginLambda = ginadapter.New(ginEngine)
	}
}

// Handler is the function that executes for every Request passed into the Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if isLocal == \"\" {
		lambda.Start(Handler)
	} else {
		ginEngine.Run(common.EnvValues.PortNumber)
	}
}" >> ./clients/internal-api.go




###### logger.go ###########

echo "package common

import (
	\"commonpkg/customLog\"
	\"log\"
	\"os\"
)

var logginObj customLog.ILogger

func InitLogger() {
	var err error
	logginObj, err = customLog.Init(1, \":$fileName: \", os.Stdout)

	if err != nil {
		log.Fatal(\"logger not initilized\")
	}
}

func WriteLog(logLevel int, logMessage string) {
	logginObj.Write(logLevel, logMessage)
}" >>./common/logging.go





###### env.go ###########
echo "package common

import \"os\"

type Env struct {
	TableName  string
	PortNumber string
}

var EnvValues Env

func InitEnv(isLocal string) {
	if isLocal == \"\" {
		EnvValues = Env{
			TableName:  os.Getenv(\"$name\"),
			PortNumber: \"0\",
		}
	} else {
		EnvValues = Env{
			TableName:  \"$fileName\",
			PortNumber: \":8080\",
		}
	}
}" >> ./common/env.go

####### controller ###########

echo "package controller

import (
	commonModels \"commonpkg/models\"
	\"net/http\"

	\"$fileName/service\"

	\"github.com/gin-gonic/gin\"
)

type AddRequest struct {
	Code string `json:\"code\" binding:\"required\"`
}

type AddMultipleRequest struct {
	Codes []string `json:\"codes\"`
}
type GetRequest struct {
	Id string `uri:\"id\"`
}

var ${privatename}Ctr *${name}Controller

type ${name}Controller struct {
	${privatename}Svc *service.${name}Service
}

func Init${name}Controller() (*${name}Controller, *commonModels.ErrorDetail) {
	if ${privatename}Ctr == nil {
		svc, err := service.Init${name}Service()
		if err != nil {
			return nil, err
		}
		${privatename}Ctr = &${name}Controller{
			${privatename}Svc: svc,
		}
	}
	return ${privatename}Ctr, nil
}
func (ctrl *${name}Controller) GetAll(context *gin.Context) {
	data := ctrl.${privatename}Svc.GetAll()
	context.JSON(data.StatusCode, data)
}

func (ctrl *${name}Controller) Get(context *gin.Context) {
	var getRquest GetRequest

	if err := context.ShouldBindUri(&getRquest); err == nil {
		data := ctrl.${privatename}Svc.Get(getRquest.Id)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			\"error\": err,
		})
	}

}

func (ctrl *${name}Controller) Add(context *gin.Context) {
	var addData AddRequest

	var b []byte
	context.Request.Body.Read(b)

	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.${privatename}Svc.Add(addData.Code)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			\"error\": err,
		})
	}
}

func (ctrl *${name}Controller) AddMultiple(context *gin.Context) {
	var addData AddMultipleRequest
	if err := context.ShouldBindJSON(&addData); err == nil {
		data := ctrl.${privatename}Svc.AddMultiple(addData.Codes)
		context.JSON(data.StatusCode, data)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			\"error\": err,
		})
	}
}" >> "./controller/$fileName.controller.go"



##### persistance ####

echo "package persistance

import (
	commonModels \"commonpkg/models\"
	\"fmt\"
	\"$fileName/common\"

	\"github.com/aws/aws-sdk-go/aws\"
	\"github.com/aws/aws-sdk-go/aws/session\"
	\"github.com/aws/aws-sdk-go/service/dynamodb\"
	\"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute\"
	uuid \"github.com/iris-contrib/go.uuid\"
)

var ${privatename}PersistanceObj *${name}Persistance

type ${name}Persistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func Init${name}Persistance() (*${name}Persistance, *commonModels.ErrorDetail) {
	if ${privatename}PersistanceObj == nil {
		dbSession, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})

		if err != nil {
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorDbConnection,
				ErrorMessage: err.Error(),
			}
		}
		dynamoDbSession := session.Must(dbSession, err)

		${privatename}PersistanceObj = &${name}Persistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return ${privatename}PersistanceObj, nil
}

func (repo *${name}Persistance) GetAll() ([]commonModels.${name}Dto, *commonModels.ErrorDetail) {

	result, err := repo.db.Scan(&dynamodb.ScanInput{
		TableName: &repo.tableName,
		Limit:     aws.Int64(100),
	})
	if err != nil {
		common.WriteLog(1, err.Error())
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: err.Error(),
		}
	}
	if result.Items == nil {
		message := \"Could not find ${name}s\"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.${name}Dto, 0)
	tempItem, errorDetails := build${name}s(result.Items)
	if errorDetails != nil {
		return nil, errorDetails
	}
	items = append(items, tempItem...)

	for result.LastEvaluatedKey != nil {
		result, err = repo.db.Scan(&dynamodb.ScanInput{
			TableName:         &repo.tableName,
			Limit:             aws.Int64(100),
			ExclusiveStartKey: result.LastEvaluatedKey,
		})
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: err.Error(),
			}
		}

		if result.Items != nil {
			tempItem, errorDetails = build${name}s(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func build${name}s(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.${name}Dto, *commonModels.ErrorDetail) {
	items := make([]commonModels.${name}Dto, 0)

	for _, val := range dbItems {
		item := commonModels.${name}Dto{}
		err := dynamodbattribute.UnmarshalMap(val, &item)
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: err.Error(),
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (repo *${name}Persistance) Get(id string) (*commonModels.${name}Dto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.${name}Dto
	result, err := repo.db.Query(&dynamodb.QueryInput{
		TableName: &repo.tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			\"id\": {
				ComparisonOperator: aws.String(\"EQ\"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: err.Error(),
		}
	}
	if result.Items == nil || len(result.Items) == 0 {
		message := fmt.Sprintf(\"Could not find ${name}s for id %s\", id)
		common.WriteLog(3, message)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}
	if len(result.Items) > 0 {

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &hnsCode)
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: \"${name} not is correct format\",
			}
		}
	}
	return hnsCode, nil
}

func (repo *${name}Persistance) Add(code string) (*commonModels.${name}Dto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	new${name} := commonModels.${name}Dto{Id: id.String(), ${name}: code}

	av, err := dynamodbattribute.MarshalMap(new${name})
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf(\"Got error marshalling new ${name} item: %s\", err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf(\"Error in adding ${name} %s, error message; %s\", code, err.Error()),
		}
	}

	return &new${name}, nil
}

func (repo *${name}Persistance) AddMultiple(codes []string) ([]commonModels.${name}Dto, []commonModels.ErrorDetail) {
	var new${name}s []commonModels.${name}Dto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		new${name}, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			new${name}s = append(new${name}s, *new${name})
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return new${name}s, nil
}" >> "./persistance/$fileName.persistance.go"

######## router #########


echo "package router

import (
	\"$fileName/common\"
	\"$fileName/controller\"

	\"github.com/gin-gonic/gin\"
)

func InitRoutes(engine *gin.Engine) {
	controller, err := controller.Init${name}Controller()
	if err != nil {
		common.WriteLog(1, err.Error())
		panic(err)
	}
	engine.GET(\"/\", func(c *gin.Context) {
		controller.GetAll(c)
	})
	engine.GET(\"/:id\", func(c *gin.Context) {
		controller.Get(c)
	})
	engine.POST(\"/\", func(c *gin.Context) {
		controller.Add(c)
	})
	engine.POST(\"/addmultiple\", func(c *gin.Context) {
		controller.AddMultiple(c)
	})
}
" >> "./router/$fileName.router.go"

##### service #######


echo "package service

import (
	commonModels \"commonpkg/models\"
	\"fmt\"
	\"net/http\"

	\"$fileName/persistance\"
)

var ${name}Obj *${name}Service

type ${name}Service struct {
	${privatename}Repo *persistance.${name}Persistance
}

func Init${name}Service() (*${name}Service, *commonModels.ErrorDetail) {
	if ${name}Obj == nil {
		repo, err := persistance.Init${name}Persistance()
		if err != nil {
			return nil, err
		}
		${name}Obj = &${name}Service{
			${privatename}Repo: repo,
		}
	}
	return ${name}Obj, nil
}

func (service *${name}Service) GetAll() commonModels.${name}ListResponse {
	allCodes, err := service.${privatename}Repo.GetAll()

	if err != nil {

		return commonModels.${name}ListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: \"could not get All ${name}\",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.${name}ListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusOK,
				},
				Start:    0,
				Total:    len(allCodes),
				PageSize: len(allCodes),
			},
			Data: allCodes,
		}
	}
}

func (service *${name}Service) Get(id string) commonModels.${name}Response {
	${privatename}, err := service.${privatename}Repo.Get(id)
	if err != nil {
		return commonModels.${name}Response{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf(\"Could not get HSN Code for id: %s\", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.${name}Response{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *${privatename},
		}
	}
}

func (service *${name}Service) Add(code string) commonModels.${name}Response {
	${privatename}, err := service.${privatename}Repo.Add(code)

	if err != nil {
		return commonModels.${name}Response{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf(\"could not add HSN Code - %s\", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.${name}Response{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *${privatename},
		}
	}
}

func (service *${name}Service) AddMultiple(codes []string) commonModels.${name}ListResponse {
	allCodes, err := service.${privatename}Repo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.${name}ListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: \"could not add All ${name}\",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.${name}ListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: \"could not add All ${name}\",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.${name}ListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				Start:    0,
				Total:    len(allCodes),
				PageSize: len(allCodes),
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusCreated,
				},
			},
			Data: allCodes,
		}
	}
}" >> "./service/$fileName.service.go"

go mod tidy

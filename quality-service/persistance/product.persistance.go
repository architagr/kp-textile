package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"quality-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	uuid "github.com/iris-contrib/go.uuid"
)

var productPersistanceObj *ProductPersistance

type ProductPersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitProductPersistance() (*ProductPersistance, *commonModels.ErrorDetail) {
	if productPersistanceObj == nil {
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

		productPersistanceObj = &ProductPersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.ProductTableName,
		}
	}

	return productPersistanceObj, nil
}

func (repo *ProductPersistance) GetAll() ([]commonModels.ProductDto, *commonModels.ErrorDetail) {
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
		message := "Could not find products"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.ProductDto, 0)
	tempItem, errorDetails := buildProduct(result.Items)
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
			tempItem, errorDetails = buildProduct(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildProduct(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.ProductDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.ProductDto, 0)

	for _, val := range dbItems {
		item := commonModels.ProductDto{}
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

func (repo *ProductPersistance) Get(id string) (*commonModels.ProductDto, *commonModels.ErrorDetail) {
	var quality *commonModels.ProductDto
	result, err := repo.db.Query(&dynamodb.QueryInput{
		TableName: &repo.tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
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
		message := fmt.Sprintf("Could not find Product for id %s", id)
		common.WriteLog(3, message)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}
	if len(result.Items) > 0 {

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &quality)
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: "Product not is correct format",
			}
		}
	}
	return quality, nil
}

func (repo *ProductPersistance) Add(code string) (*commonModels.ProductDto, *commonModels.ErrorDetail) {

	filter := expression.Name("name").Equal(expression.Value(code))
	builder := expression.NewBuilder().WithFilter(filter)
	expr, err := builder.Build()

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("Got error building expression: %s", err.Error()))
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInvalidRequestParam,
			ErrorMessage: "Error building filter",
		}
	}
	result, _ := repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.tableName),
	})

	if len(result.Items) != 0 {
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Product %s already exists", code),
		}
	}

	id, _ := uuid.NewV1()
	newQualityService := commonModels.ProductDto{Id: id.String(), Name: code}

	av, err := dynamodbattribute.MarshalMap(newQualityService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new Product item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding Product %s, error message; %s", code, err.Error()),
		}
	}

	return &newQualityService, nil
}

func (repo *ProductPersistance) AddMultiple(codes []string) ([]commonModels.ProductDto, []commonModels.ErrorDetail) {
	var newProducts []commonModels.ProductDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newProductService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newProducts = append(newProducts, *newProductService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newProducts, nil
}

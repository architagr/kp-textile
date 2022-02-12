package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"organization-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	uuid "github.com/iris-contrib/go.uuid"
)

var godownPersistanceObj *GodownPersistance

type GodownPersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitGodownPersistance() (*GodownPersistance, *commonModels.ErrorDetail) {
	if godownPersistanceObj == nil {
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

		godownPersistanceObj = &GodownPersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.GodownTableName,
		}
	}

	return godownPersistanceObj, nil
}

func (repo *GodownPersistance) GetAll() ([]commonModels.GodownDto, *commonModels.ErrorDetail) {
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
		message := "Could not find Godown"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.GodownDto, 0)
	tempItem, errorDetails := buildGodown(result.Items)
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
			tempItem, errorDetails = buildGodown(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildGodown(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.GodownDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.GodownDto, 0)

	for _, val := range dbItems {
		item := commonModels.GodownDto{}
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

func (repo *GodownPersistance) Add(name string) (*commonModels.GodownDto, *commonModels.ErrorDetail) {
	filter := expression.Name("name").Equal(expression.Value(name))
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
			ErrorMessage: fmt.Sprintf("Godown: %s already exists", name),
		}
	}

	id, _ := uuid.NewV1()
	newGodown := commonModels.GodownDto{Id: id.String(), Name: name}

	av, err := dynamodbattribute.MarshalMap(newGodown)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new godown: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding godown %s, error message; %s", name, err.Error()),
		}
	}

	return &newGodown, nil
}

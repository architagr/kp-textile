package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var itemServicePersistanceObj *ItemServicePersistance

type ItemServicePersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitItemServicePersistance() (*ItemServicePersistance, *commonModels.ErrorDetail) {
	if itemServicePersistanceObj == nil {
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

		itemServicePersistanceObj = &ItemServicePersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return itemServicePersistanceObj, nil
}

func (repo *ItemServicePersistance) GetAll() ([]commonModels.ItemServiceDto, *commonModels.ErrorDetail) {

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
		message := "Could not find ItemServices"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.ItemServiceDto, 0)
	tempItem, errorDetails := buildItemServices(result.Items)
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
			tempItem, errorDetails = buildItemServices(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildItemServices(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.ItemServiceDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.ItemServiceDto, 0)

	for _, val := range dbItems {
		item := commonModels.ItemServiceDto{}
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

func (repo *ItemServicePersistance) Get(id string) (*commonModels.ItemServiceDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.ItemServiceDto
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
		message := fmt.Sprintf("Could not find ItemServices for id %s", id)
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
				ErrorMessage: "ItemService not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *ItemServicePersistance) Add(code string) (*commonModels.ItemServiceDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newItemService := commonModels.ItemServiceDto{Id: id.String(), ItemService: code}

	av, err := dynamodbattribute.MarshalMap(newItemService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new ItemService item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding ItemService %s, error message; %s", code, err.Error()),
		}
	}

	return &newItemService, nil
}

func (repo *ItemServicePersistance) AddMultiple(codes []string) ([]commonModels.ItemServiceDto, []commonModels.ErrorDetail) {
	var newItemServices []commonModels.ItemServiceDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newItemService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newItemServices = append(newItemServices, *newItemService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newItemServices, nil
}

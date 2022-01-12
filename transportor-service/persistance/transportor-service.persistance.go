package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"transportor-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var transportorServicePersistanceObj *TransportorServicePersistance

type TransportorServicePersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitTransportorServicePersistance() (*TransportorServicePersistance, *commonModels.ErrorDetail) {
	if transportorServicePersistanceObj == nil {
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

		transportorServicePersistanceObj = &TransportorServicePersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return transportorServicePersistanceObj, nil
}

func (repo *TransportorServicePersistance) GetAll() ([]commonModels.TransportorServiceDto, *commonModels.ErrorDetail) {

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
		message := "Could not find TransportorServices"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.TransportorServiceDto, 0)
	tempItem, errorDetails := buildTransportorServices(result.Items)
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
			tempItem, errorDetails = buildTransportorServices(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildTransportorServices(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.TransportorServiceDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.TransportorServiceDto, 0)

	for _, val := range dbItems {
		item := commonModels.TransportorServiceDto{}
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

func (repo *TransportorServicePersistance) Get(id string) (*commonModels.TransportorServiceDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.TransportorServiceDto
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
		message := fmt.Sprintf("Could not find TransportorServices for id %s", id)
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
				ErrorMessage: "TransportorService not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *TransportorServicePersistance) Add(code string) (*commonModels.TransportorServiceDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newTransportorService := commonModels.TransportorServiceDto{Id: id.String(), TransportorService: code}

	av, err := dynamodbattribute.MarshalMap(newTransportorService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new TransportorService item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding TransportorService %s, error message; %s", code, err.Error()),
		}
	}

	return &newTransportorService, nil
}

func (repo *TransportorServicePersistance) AddMultiple(codes []string) ([]commonModels.TransportorServiceDto, []commonModels.ErrorDetail) {
	var newTransportorServices []commonModels.TransportorServiceDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newTransportorService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newTransportorServices = append(newTransportorServices, *newTransportorService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newTransportorServices, nil
}

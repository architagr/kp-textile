package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"document-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var documentServicePersistanceObj *DocumentServicePersistance

type DocumentServicePersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitDocumentServicePersistance() (*DocumentServicePersistance, *commonModels.ErrorDetail) {
	if documentServicePersistanceObj == nil {
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

		documentServicePersistanceObj = &DocumentServicePersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return documentServicePersistanceObj, nil
}

func (repo *DocumentServicePersistance) GetAll() ([]commonModels.DocumentServiceDto, *commonModels.ErrorDetail) {

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
		message := "Could not find DocumentServices"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.DocumentServiceDto, 0)
	tempItem, errorDetails := buildDocumentServices(result.Items)
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
			tempItem, errorDetails = buildDocumentServices(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildDocumentServices(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.DocumentServiceDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.DocumentServiceDto, 0)

	for _, val := range dbItems {
		item := commonModels.DocumentServiceDto{}
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

func (repo *DocumentServicePersistance) Get(id string) (*commonModels.DocumentServiceDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.DocumentServiceDto
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
		message := fmt.Sprintf("Could not find DocumentServices for id %s", id)
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
				ErrorMessage: "DocumentService not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *DocumentServicePersistance) Add(code string) (*commonModels.DocumentServiceDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newDocumentService := commonModels.DocumentServiceDto{Id: id.String(), DocumentService: code}

	av, err := dynamodbattribute.MarshalMap(newDocumentService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new DocumentService item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding DocumentService %s, error message; %s", code, err.Error()),
		}
	}

	return &newDocumentService, nil
}

func (repo *DocumentServicePersistance) AddMultiple(codes []string) ([]commonModels.DocumentServiceDto, []commonModels.ErrorDetail) {
	var newDocumentServices []commonModels.DocumentServiceDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newDocumentService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newDocumentServices = append(newDocumentServices, *newDocumentService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newDocumentServices, nil
}

package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"hsn-code-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var hnsCodepersistanceObj *HnsCodePersistance

type HnsCodePersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitHnsCodePersistance() (*HnsCodePersistance, *commonModels.ErrorDetail) {
	if hnsCodepersistanceObj == nil {
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

		hnsCodepersistanceObj = &HnsCodePersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return hnsCodepersistanceObj, nil
}

func (repo *HnsCodePersistance) GetAll() ([]commonModels.HnsCodeDto, *commonModels.ErrorDetail) {

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
		message := "Could not find HSN codes"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.HnsCodeDto, 0)
	tempItem, errorDetails := buildHsnCodes(result.Items)
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
			tempItem, errorDetails = buildHsnCodes(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildHsnCodes(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.HnsCodeDto, 0)

	for _, val := range dbItems {
		item := commonModels.HnsCodeDto{}
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

func (repo *HnsCodePersistance) Get(id string) (*commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.HnsCodeDto
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
		message := fmt.Sprintf("Could not find HSN codes for id %s", id)
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
				ErrorMessage: "hsn codes not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *HnsCodePersistance) Add(code string) (*commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newHnsCode := commonModels.HnsCodeDto{Id: id.String(), HnsCode: code}

	av, err := dynamodbattribute.MarshalMap(newHnsCode)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new HnsCode item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding HSN code %s, error message; %s", code, err.Error()),
		}
	}

	return &newHnsCode, nil
}

func (repo *HnsCodePersistance) AddMultiple(codes []string) ([]commonModels.HnsCodeDto, []commonModels.ErrorDetail) {
	var newHnsCodes []commonModels.HnsCodeDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newHnsCode, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newHnsCodes = append(newHnsCodes, *newHnsCode)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newHnsCodes, nil
}

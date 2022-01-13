package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"quality-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var qualityPersistanceObj *QualityPersistance

type QualityPersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitQualityPersistance() (*QualityPersistance, *commonModels.ErrorDetail) {
	if qualityPersistanceObj == nil {
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

		qualityPersistanceObj = &QualityPersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return qualityPersistanceObj, nil
}

func (repo *QualityPersistance) GetAll() ([]commonModels.QualityDto, *commonModels.ErrorDetail) {

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
		message := "Could not find Qualities"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.QualityDto, 0)
	tempItem, errorDetails := buildQuality(result.Items)
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
			tempItem, errorDetails = buildQuality(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildQuality(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.QualityDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.QualityDto, 0)

	for _, val := range dbItems {
		item := commonModels.QualityDto{}
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

func (repo *QualityPersistance) Get(id string) (*commonModels.QualityDto, *commonModels.ErrorDetail) {
	var quality *commonModels.QualityDto
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
		message := fmt.Sprintf("Could not find Quality for id %s", id)
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
				ErrorMessage: "Quality not is correct format",
			}
		}
	}
	return quality, nil
}

func (repo *QualityPersistance) Add(code string) (*commonModels.QualityDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newQualityService := commonModels.QualityDto{Id: id.String(), Name: code}

	av, err := dynamodbattribute.MarshalMap(newQualityService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new Quality item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding Quality %s, error message; %s", code, err.Error()),
		}
	}

	return &newQualityService, nil
}

func (repo *QualityPersistance) AddMultiple(codes []string) ([]commonModels.QualityDto, []commonModels.ErrorDetail) {
	var newQualityServices []commonModels.QualityDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newQualityService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newQualityServices = append(newQualityServices, *newQualityService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newQualityServices, nil
}

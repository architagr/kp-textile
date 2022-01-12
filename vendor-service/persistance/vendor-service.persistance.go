package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"vendor-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var vendorServicePersistanceObj *VendorServicePersistance

type VendorServicePersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitVendorServicePersistance() (*VendorServicePersistance, *commonModels.ErrorDetail) {
	if vendorServicePersistanceObj == nil {
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

		vendorServicePersistanceObj = &VendorServicePersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.TableName,
		}
	}

	return vendorServicePersistanceObj, nil
}

func (repo *VendorServicePersistance) GetAll() ([]commonModels.VendorServiceDto, *commonModels.ErrorDetail) {

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
		message := "Could not find VendorServices"
		common.WriteLog(5, message)

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}

	items := make([]commonModels.VendorServiceDto, 0)
	tempItem, errorDetails := buildVendorServices(result.Items)
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
			tempItem, errorDetails = buildVendorServices(result.Items)
			if errorDetails != nil {
				return nil, errorDetails
			}
			items = append(items, tempItem...)
		}
	}
	return items, nil
}

func buildVendorServices(dbItems []map[string]*dynamodb.AttributeValue) ([]commonModels.VendorServiceDto, *commonModels.ErrorDetail) {
	items := make([]commonModels.VendorServiceDto, 0)

	for _, val := range dbItems {
		item := commonModels.VendorServiceDto{}
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

func (repo *VendorServicePersistance) Get(id string) (*commonModels.VendorServiceDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.VendorServiceDto
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
		message := fmt.Sprintf("Could not find VendorServices for id %s", id)
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
				ErrorMessage: "VendorService not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *VendorServicePersistance) Add(code string) (*commonModels.VendorServiceDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newVendorService := commonModels.VendorServiceDto{Id: id.String(), VendorService: code}

	av, err := dynamodbattribute.MarshalMap(newVendorService)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new VendorService item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding VendorService %s, error message; %s", code, err.Error()),
		}
	}

	return &newVendorService, nil
}

func (repo *VendorServicePersistance) AddMultiple(codes []string) ([]commonModels.VendorServiceDto, []commonModels.ErrorDetail) {
	var newVendorServices []commonModels.VendorServiceDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newVendorService, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newVendorServices = append(newVendorServices, *newVendorService)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newVendorServices, nil
}

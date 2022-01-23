package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var purchansePersistanceObj *PurchasePersistance

type PurchasePersistance struct {
	db                 *dynamodb.DynamoDB
	inventoryTableName string
	bailInfoTable      string
	itemTable          string
}

func InitPurchasePersistance() (*PurchasePersistance, *commonModels.ErrorDetail) {
	if purchansePersistanceObj == nil {
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

		purchansePersistanceObj = &PurchasePersistance{
			db:                 dynamodb.New(dynamoDbSession),
			inventoryTableName: common.EnvValues.InventoryTableName,
			bailInfoTable:      common.EnvValues.BailInfoTableName,
			itemTable:          common.EnvValues.ItemTableName,
		}
	}

	return purchansePersistanceObj, nil
}

func (repo *PurchasePersistance) GetAllPurchaseOrders(request commonModels.InventoryListRequest) ([]commonModels.InventoryDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey("")),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getInventoryDetailError := getInventoryDetails(expr, request)
	if getInventoryDetailError != nil {
		return nil, nil, getInventoryDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	inventoryDetails, inventoryListParseErr := parseDbItemsToInventoryList(result.Items)
	if inventoryListParseErr != nil {
		return nil, nil, inventoryListParseErr
	}
	for len(inventoryDetails) < int(request.PageSize) && request.LastEvalutionKey != nil {
		result, getInventoryDetailError = getInventoryDetails(expr, request)
		if getInventoryDetailError != nil {
			return nil, nil, getInventoryDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		inventoryDetailsTemp, inventoryListParseErr := parseDbItemsToInventoryList(result.Items)
		if inventoryListParseErr != nil {
			return nil, nil, inventoryListParseErr
		}
		inventoryDetails = append(inventoryDetails, inventoryDetailsTemp...)
	}

	return inventoryDetails, request.LastEvalutionKey, nil
}

func (repo *PurchasePersistance) GetTotalPurchaseOrders(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
	var count int64 = 0
	proj := expression.NamesList(expression.Name("branchId"))
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey("")),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(proj).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return 0, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	request.PageSize = 100
	result, getInventoryDetailError := getInventoryDetails(expr, request)
	if getInventoryDetailError != nil {
		return 0, getInventoryDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))

	for len(result.Items) > 0 && request.LastEvalutionKey != nil {
		result, getInventoryDetailError = getInventoryDetails(expr, request)
		if getInventoryDetailError != nil {
			return 0, getInventoryDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}

	return count, nil
}

func (repo *PurchasePersistance) GetPurchaseBillDetails(request commonModels.InventoryFilterDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey(request.PurchaseBillNumber)),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getInventoryDetailError := getInventoryDetails(expr, commonModels.InventoryListRequest{
		InventoryFilterDto: request,
	})

	if getInventoryDetailError != nil {
		return nil, getInventoryDetailError
	}

	if len(result.Items) > 0 {
		inventory, err := parseDbItemToInventory(result.Items[0])
		if err != nil {
			return nil, err
		}
		return inventory, nil
	}

	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: fmt.Sprintf("No purchase order found by the order no %s", request.PurchaseBillNumber),
	}
}

func (repo *PurchasePersistance) UpsertPurchaseOrder(data commonModels.InventoryDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling purchase details, purchase bill number - %s, branch id - %s, err: %s", data.BillNo, data.BranchId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.inventoryTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating purchase data for bill no %s, branch id %s, error message; %s", data.BillNo, data.BranchId, err.Error()),
		}
	}
	return &data, nil
}

func (repo *PurchasePersistance) DeletePurchaseBillDetails(branchId, billno string) *commonModels.ErrorDetail {
	purchanseDetails, getPurchaseDetailErr := repo.GetPurchaseBillDetails(commonModels.InventoryFilterDto{
		BranchId:           branchId,
		PurchaseBillNumber: billno,
	})
	if getPurchaseDetailErr != nil {
		return getPurchaseDetailErr
	}

	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.inventoryTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"inventorySortKey": {
				S: aws.String(purchanseDetails.InventorySortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting bill no: %s for branchId: %s, error; %s", billno, branchId, err.Error()),
		}
	}
	timeNow := time.Now().UTC().Unix()
	purchanseDetails.InventorySortKey = fmt.Sprintf("%s|%d", common.GetInventoryDeleteSortKey(billno), timeNow)
	_, updateError := repo.UpsertPurchaseOrder(*purchanseDetails)
	if updateError != nil {
		return updateError
	}

	return nil
}

// #region private functions
func getInventoryDetails(expr expression.Expression, request commonModels.InventoryListRequest) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExclusiveStartKey:         request.LastEvalutionKey,
		TableName:                 aws.String(purchansePersistanceObj.inventoryTableName),
	}

	if request.PageSize > 0 {
		queryInput.Limit = &request.PageSize
	}
	result, err := purchansePersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get purchase bill for branch %s call failed: %s", request.BranchId, err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}

func parseDbItemsToInventoryList(items []map[string]*dynamodb.AttributeValue) ([]commonModels.InventoryDto, *commonModels.ErrorDetail) {
	inventoryDetails := make([]commonModels.InventoryDto, 0)

	if len(items) > 0 {
		for _, val := range items {
			inventory, err := parseDbItemToInventory(val)
			if err != nil {
				return nil, err
			}
			inventoryDetails = append(inventoryDetails, *inventory)
		}
		return inventoryDetails, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: "No data found",
	}
}
func parseDbItemToInventory(item map[string]*dynamodb.AttributeValue) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
	inventory := commonModels.InventoryDto{}

	err := dynamodbattribute.UnmarshalMap(item, &inventory)

	if err != nil {
		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}

	return &inventory, nil
}

// #endregion

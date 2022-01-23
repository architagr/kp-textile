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

var salesPersistanceObj *SalesPersistance

type SalesPersistance struct {
	db                 *dynamodb.DynamoDB
	inventoryTableName string
	bailInfoTable      string
	itemTable          string
}

func InitSalesPersistance() (*SalesPersistance, *commonModels.ErrorDetail) {
	if salesPersistanceObj == nil {
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

		salesPersistanceObj = &SalesPersistance{
			db:                 dynamodb.New(dynamoDbSession),
			inventoryTableName: common.EnvValues.InventoryTableName,
			bailInfoTable:      common.EnvValues.BailInfoTableName,
			itemTable:          common.EnvValues.ItemTableName,
		}
	}

	return salesPersistanceObj, nil
}

func (repo *SalesPersistance) GetAllSalesOrders(request commonModels.InventoryListRequest) ([]commonModels.InventoryDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey("")),
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

	return inventoryDetails, result.LastEvaluatedKey, nil
}

func (repo *SalesPersistance) GetTotalSalesOrders(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
	var count int64 = 0
	proj := expression.NamesList(expression.Name("branchId"))
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey("")),
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

func (repo *SalesPersistance) GetSalesBillDetails(request commonModels.InventoryFilterDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey(request.SalesBillNumber)),
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
		ErrorMessage: fmt.Sprintf("No Sales order found by the order no %s", request.SalesBillNumber),
	}
}

func (repo *SalesPersistance) UpsertSalesOrder(data commonModels.InventoryDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling purchase details, sales bill number - %s, branch id - %s, err: %s", data.BillNo, data.BranchId, err),
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
			ErrorMessage: fmt.Sprintf("Error in adding/updating sales data for bill no %s, branch id %s, error message; %s", data.BillNo, data.BranchId, err.Error()),
		}
	}
	return &data, nil
}

func (repo *SalesPersistance) DeleteSalesBillDetails(branchId, billno string) *commonModels.ErrorDetail {
	salesDetails, getSalesDetailErr := repo.GetSalesBillDetails(commonModels.InventoryFilterDto{
		BranchId:        branchId,
		SalesBillNumber: billno,
	})
	if getSalesDetailErr != nil {
		return getSalesDetailErr
	}

	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.inventoryTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"inventorySortKey": {
				S: aws.String(salesDetails.InventorySortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting sales bill no: %s for branchId: %s, error: %s", billno, branchId, err.Error()),
		}
	}

	timeNow := time.Now().UTC().Unix()
	salesDetails.InventorySortKey = fmt.Sprintf("%s|%s|%d", common.GetInventoryDeleteSortKey(billno), common.SORTKEY_INVENTORY_SALES, timeNow)
	_, updateError := repo.UpsertSalesOrder(*salesDetails)
	if updateError != nil {
		return updateError
	}

	return nil
}

package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ISalesPersistance interface {
	GetAllTotal(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail)
	GetAll(request commonModels.InventoryListRequest) ([]commonModels.SalesDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail)
	GetById(salesId string) (*commonModels.SalesDto, *commonModels.ErrorDetail)
	GetByBillNo(salesBillNo string) (*commonModels.SalesDto, *commonModels.ErrorDetail)
	Add(data commonModels.SalesDto) (*commonModels.SalesDto, *commonModels.ErrorDetail)
}

var salesPersistanceObj *SalesPersistance

type SalesPersistance struct {
	db                   *dynamodb.DynamoDB
	salesTableName       string
	salesIdIndexName     string
	salesBillNoIndexName string
	challanNoIndexName   string
}

func InitSalesPersistance() (ISalesPersistance, *commonModels.ErrorDetail) {
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
			db:                   dynamodb.New(dynamoDbSession),
			salesTableName:       common.EnvValues.SalesTableName,
			salesIdIndexName:     common.EnvValues.SalesIdIndexName,
			salesBillNoIndexName: common.EnvValues.SalesBillNoIndexName,
			challanNoIndexName:   common.EnvValues.ChallanNoIndexName,
		}
	}

	return salesPersistanceObj, nil
}

func (repo *SalesPersistance) GetAllTotal(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
	var count int64 = 0
	proj := expression.NamesList(expression.Name("godownId"))
	keyCondition := expression.Key("godownId").Equal(expression.Value(request.GodownId))
	if len(request.ProductId) > 0 {
		sortKey := common.GetSalesSortKey(request.ProductId, request.QualityId, "")
		keyCondition.And(expression.Key("sortKey").BeginsWith(sortKey))
	}
	exprBuilder := expression.NewBuilder()

	if len(request.SalesBillNumber) > 0 {
		filter := expression.Name("salesBillNo").Equal(expression.Value(request.SalesBillNumber))
		exprBuilder = exprBuilder.WithFilter(filter)
	}

	expr, err := exprBuilder.WithKeyCondition(keyCondition).WithProjection(proj).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return 0, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	request.PageSize = 100
	result, getInventoryDetailError := getSalesDetails(expr, request)
	if getInventoryDetailError != nil {
		return 0, getInventoryDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))

	for len(result.Items) > 0 && request.LastEvalutionKey != nil {
		result, getInventoryDetailError = getSalesDetails(expr, request)
		if getInventoryDetailError != nil {
			return 0, getInventoryDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}

	return count, nil
}

func (repo *SalesPersistance) GetAll(request commonModels.InventoryListRequest) ([]commonModels.SalesDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {

	keyCondition := expression.Key("godownId").Equal(expression.Value(request.GodownId))
	if len(request.ProductId) > 0 {
		sortKey := common.GetSalesSortKey(request.ProductId, request.QualityId, "")
		keyCondition.And(expression.Key("sortKey").BeginsWith(sortKey))
	}

	exprBuilder := expression.NewBuilder()

	if len(request.SalesBillNumber) > 0 {
		filter := expression.Name("salesBillNo").Equal(expression.Value(request.SalesBillNumber))
		exprBuilder = exprBuilder.WithFilter(filter)
	}
	expr, err := exprBuilder.WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getSalesDetailError := getSalesDetails(expr, request)
	if getSalesDetailError != nil {
		return nil, nil, getSalesDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	salesDetails, salesListParseErr := parseDbItemsToSalesList(result.Items)
	if salesListParseErr != nil {
		return nil, nil, salesListParseErr
	}
	for len(salesDetails) < int(request.PageSize) && request.LastEvalutionKey != nil {
		result, getSalesDetailError = getSalesDetails(expr, request)
		if getSalesDetailError != nil {
			return nil, nil, getSalesDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		inventoryDetailsTemp, inventoryListParseErr := parseDbItemsToSalesList(result.Items)
		if inventoryListParseErr != nil {
			return nil, nil, inventoryListParseErr
		}
		salesDetails = append(salesDetails, inventoryDetailsTemp...)
	}
	return salesDetails, request.LastEvalutionKey, nil
}

func (repo *SalesPersistance) GetById(salesId string) (*commonModels.SalesDto, *commonModels.ErrorDetail) {
	keyCondition := expression.Key("salesId").Equal(expression.Value(salesId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getSalesDetailError := getSaleDetailsInIndex(expr, commonModels.InventoryListRequest{
		SalesId:  salesId,
		PageSize: 0,
	}, repo.salesIdIndexName)

	if getSalesDetailError != nil {
		return nil, getSalesDetailError
	}

	salesDetails, salesListParseErr := parseDbItemsToSalesList(result.Items)
	if salesListParseErr != nil {
		return nil, salesListParseErr
	}

	return &salesDetails[0], nil
}
func (repo *SalesPersistance) GetByBillNo(salesBillNo string) (*commonModels.SalesDto, *commonModels.ErrorDetail) {
	keyCondition := expression.Key("salesBillNo").Equal(expression.Value(salesBillNo))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getSalesDetailError := getSaleDetailsInIndex(expr, commonModels.InventoryListRequest{
		InventoryFilterDto: commonModels.InventoryFilterDto{
			SalesBillNumber: salesBillNo,
		},
		PageSize: 0,
	}, repo.salesBillNoIndexName)

	if getSalesDetailError != nil {
		return nil, getSalesDetailError
	}

	salesDetails, salesListParseErr := parseDbItemsToSalesList(result.Items)
	if salesListParseErr != nil {
		return nil, salesListParseErr
	}

	return &salesDetails[0], nil
}

func (repo *SalesPersistance) Add(data commonModels.SalesDto) (*commonModels.SalesDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling sales details, sales bill number - %s, godown id - %s, err: %s", data.SalesBillNo, data.GodownId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.salesTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating sales data for bill no %s, godown id %s, error message; %s", data.SalesBillNo, data.GodownId, err.Error()),
		}
	}
	return &data, nil
}

func getSalesDetails(expr expression.Expression, request commonModels.InventoryListRequest) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExclusiveStartKey:         request.LastEvalutionKey,
		TableName:                 aws.String(salesPersistanceObj.salesTableName),
	}

	if request.PageSize > 0 {
		queryInput.Limit = &request.PageSize
	}
	result, err := salesPersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get sales bill for godown %s call failed: %s", request.GodownId, err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}

func getSaleDetailsInIndex(expr expression.Expression, request commonModels.InventoryListRequest, indexNane string) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 &indexNane,
		TableName:                 aws.String(salesPersistanceObj.salesTableName),
	}

	if request.PageSize > 0 {
		queryInput.Limit = &request.PageSize
	}
	result, err := salesPersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get sales bill for %+v call failed: %s", request, err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}
func parseDbItemsToSalesList(items []map[string]*dynamodb.AttributeValue) ([]commonModels.SalesDto, *commonModels.ErrorDetail) {
	saleDetails := make([]commonModels.SalesDto, 0)

	if len(items) > 0 {
		for _, val := range items {
			sale, err := parseDbItemToSales(val)
			if err != nil {
				return nil, err
			}
			saleDetails = append(saleDetails, *sale)
		}
		return saleDetails, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: "No data found",
	}
}
func parseDbItemToSales(item map[string]*dynamodb.AttributeValue) (*commonModels.SalesDto, *commonModels.ErrorDetail) {
	sale := commonModels.SalesDto{}

	err := dynamodbattribute.UnmarshalMap(item, &sale)

	if err != nil {
		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}

	return &sale, nil
}

// func (repo *SalesPersistance) GetAllSalesOrders(request commonModels.InventoryListRequest) ([]commonModels.InventoryDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey("")),
// 	)

// 	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return nil, nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	result, getInventoryDetailError := getInventoryDetails(expr, request)
// 	if getInventoryDetailError != nil {
// 		return nil, nil, getInventoryDetailError
// 	}

// 	request.LastEvalutionKey = result.LastEvaluatedKey
// 	inventoryDetails, inventoryListParseErr := parseDbItemsToPurchaseList(result.Items)
// 	if inventoryListParseErr != nil {
// 		return nil, nil, inventoryListParseErr
// 	}
// 	for len(inventoryDetails) < int(request.PageSize) && request.LastEvalutionKey != nil {
// 		result, getInventoryDetailError = getInventoryDetails(expr, request)
// 		if getInventoryDetailError != nil {
// 			return nil, nil, getInventoryDetailError
// 		}
// 		request.LastEvalutionKey = result.LastEvaluatedKey
// 		inventoryDetailsTemp, inventoryListParseErr := parseDbItemsToPurchaseList(result.Items)
// 		if inventoryListParseErr != nil {
// 			return nil, nil, inventoryListParseErr
// 		}
// 		inventoryDetails = append(inventoryDetails, inventoryDetailsTemp...)

// 	}

// 	return inventoryDetails, result.LastEvaluatedKey, nil
// }

// func (repo *SalesPersistance) GetTotalSalesOrders(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
// 	var count int64 = 0
// 	proj := expression.NamesList(expression.Name("godownId"))
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey("")),
// 	)

// 	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(proj).Build()

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return 0, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	request.PageSize = 100
// 	result, getInventoryDetailError := getInventoryDetails(expr, request)
// 	if getInventoryDetailError != nil {
// 		return 0, getInventoryDetailError
// 	}

// 	request.LastEvalutionKey = result.LastEvaluatedKey
// 	count = count + int64(len(result.Items))

// 	for len(result.Items) > 0 && request.LastEvalutionKey != nil {
// 		result, getInventoryDetailError = getInventoryDetails(expr, request)
// 		if getInventoryDetailError != nil {
// 			return 0, getInventoryDetailError
// 		}
// 		request.LastEvalutionKey = result.LastEvaluatedKey
// 		count = count + int64(len(result.Items))
// 	}

// 	return count, nil
// }

// func (repo *SalesPersistance) GetSalesBillDetails(request commonModels.InventoryFilterDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventorySalesSortKey(request.SalesBillNumber)),
// 	)

// 	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	result, getInventoryDetailError := getInventoryDetails(expr, commonModels.InventoryListRequest{
// 		InventoryFilterDto: request,
// 	})

// 	if getInventoryDetailError != nil {
// 		return nil, getInventoryDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		inventory, err := parseDbItemToPurchase(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return inventory, nil
// 	}

// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("No Sales order found by the order no %s", request.SalesBillNumber),
// 	}
// }

// func (repo *SalesPersistance) GetDeletedSalesBillDetails(request commonModels.InventoryFilterDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {

// 	var inventorySortKey = fmt.Sprintf("%s|%s|", common.GetInventoryDeleteSortKey(request.SalesBillNumber), common.SORTKEY_INVENTORY_SALES)
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(inventorySortKey),
// 	)

// 	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	result, getInventoryDetailError := getInventoryDetails(expr, commonModels.InventoryListRequest{
// 		InventoryFilterDto: request,
// 	})

// 	if getInventoryDetailError != nil {
// 		return nil, getInventoryDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		inventory, err := parseDbItemToPurchase(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return inventory, nil
// 	}

// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("No deleted Sales order found by the order no %s", request.SalesBillNumber),
// 	}
// }

// func (repo *SalesPersistance) UpsertSalesOrder(data commonModels.InventoryDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
// 	av, err := dynamodbattribute.MarshalMap(data)
// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("Got error marshalling purchase details, sales bill number - %s, branch id - %s, err: %s", data.BillNo, data.GodownId, err),
// 		}
// 	}
// 	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
// 		TableName: &repo.inventoryTableName,
// 		Item:      av,
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorInsert,
// 			ErrorMessage: fmt.Sprintf("Error in adding/updating sales data for bill no %s, branch id %s, error message; %s", data.BillNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	return &data, nil
// }

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

var purchansePersistanceObj *PurchasePersistance

type PurchasePersistance struct {
	db                      *dynamodb.DynamoDB
	purchaseTableName       string
	purchaseIdIndexName     string
	purchaseBillNoIndexName string
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
			db:                      dynamodb.New(dynamoDbSession),
			purchaseTableName:       common.EnvValues.PurchaseTableName,
			purchaseIdIndexName:     common.EnvValues.PurchaseIdIndexName,
			purchaseBillNoIndexName: common.EnvValues.PurchaseBillNoIndexName,
		}
	}

	return purchansePersistanceObj, nil
}

func (repo *PurchasePersistance) GetAllTotal(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
	var count int64 = 0
	proj := expression.NamesList(expression.Name("godownId"))
	keyCondition := expression.Key("godownId").Equal(expression.Value(request.GodownId))
	if len(request.ProductId) > 0 {
		sortKey := common.GetPurchaseSortKey(request.ProductId, request.QualityId)
		keyCondition.And(expression.Key("sortKey").BeginsWith(sortKey))
	}

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
	result, getInventoryDetailError := getPurchaseDetails(expr, request)
	if getInventoryDetailError != nil {
		return 0, getInventoryDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))

	for len(result.Items) > 0 && request.LastEvalutionKey != nil {
		result, getInventoryDetailError = getPurchaseDetails(expr, request)
		if getInventoryDetailError != nil {
			return 0, getInventoryDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}

	return count, nil
}

func (repo *PurchasePersistance) GetAll(request commonModels.InventoryListRequest) ([]commonModels.PurchaseDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {
	keyCondition := expression.Key("godownId").Equal(expression.Value(request.GodownId))
	if len(request.ProductId) > 0 {
		sortKey := common.GetPurchaseSortKey(request.ProductId, request.QualityId)
		keyCondition.And(expression.Key("sortKey").BeginsWith(sortKey))
	}
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getPurchaseDetailError := getPurchaseDetails(expr, request)
	if getPurchaseDetailError != nil {
		return nil, nil, getPurchaseDetailError
	}

	request.LastEvalutionKey = result.LastEvaluatedKey
	purchaseDetails, purchaseListParseErr := parseDbItemsToPurchaseList(result.Items)
	if purchaseListParseErr != nil {
		return nil, nil, purchaseListParseErr
	}
	for len(purchaseDetails) < int(request.PageSize) && request.LastEvalutionKey != nil {
		result, getPurchaseDetailError = getPurchaseDetails(expr, request)
		if getPurchaseDetailError != nil {
			return nil, nil, getPurchaseDetailError
		}
		request.LastEvalutionKey = result.LastEvaluatedKey
		inventoryDetailsTemp, inventoryListParseErr := parseDbItemsToPurchaseList(result.Items)
		if inventoryListParseErr != nil {
			return nil, nil, inventoryListParseErr
		}
		purchaseDetails = append(purchaseDetails, inventoryDetailsTemp...)
	}
	return purchaseDetails, request.LastEvalutionKey, nil
}

func (repo *PurchasePersistance) GetById(purchaseId string) (*commonModels.PurchaseDto, *commonModels.ErrorDetail) {
	keyCondition := expression.Key("purchaseId").Equal(expression.Value(purchaseId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		errMessage := fmt.Sprintf("Got error building expression: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	result, getPurchaseDetailError := getPurchaseDetailsInIndex(expr, commonModels.InventoryListRequest{
		PurchaseId: purchaseId,
		PageSize:   0,
	}, repo.purchaseIdIndexName)

	if getPurchaseDetailError != nil {
		return nil, getPurchaseDetailError
	}

	purchaseDetails, purchaseListParseErr := parseDbItemsToPurchaseList(result.Items)
	if purchaseListParseErr != nil {
		return nil, purchaseListParseErr
	}

	return &purchaseDetails[0], nil
}

func (repo *PurchasePersistance) Add(data commonModels.PurchaseDto) (*commonModels.PurchaseDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling purchase details, purchase bill number - %s, godown id - %s, err: %s", data.PurchaseBillNo, data.GodownId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.purchaseTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating purchase data for bill no %s, godown id %s, error message; %s", data.PurchaseBillNo, data.GodownId, err.Error()),
		}
	}
	return &data, nil
}

// func (repo *PurchasePersistance) GetAllPurchaseOrders(request commonModels.InventoryListRequest) ([]commonModels.InventoryDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey("")),
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
// 	result, getInventoryDetailError := getPurchaseDetails(expr, request)
// 	if getInventoryDetailError != nil {
// 		return nil, nil, getInventoryDetailError
// 	}

// 	request.LastEvalutionKey = result.LastEvaluatedKey
// 	inventoryDetails, inventoryListParseErr := parseDbItemsToInventoryList(result.Items)
// 	if inventoryListParseErr != nil {
// 		return nil, nil, inventoryListParseErr
// 	}
// 	for len(inventoryDetails) < int(request.PageSize) && request.LastEvalutionKey != nil {
// 		result, getInventoryDetailError = getPurchaseDetails(expr, request)
// 		if getInventoryDetailError != nil {
// 			return nil, nil, getInventoryDetailError
// 		}
// 		request.LastEvalutionKey = result.LastEvaluatedKey
// 		inventoryDetailsTemp, inventoryListParseErr := parseDbItemsToInventoryList(result.Items)
// 		if inventoryListParseErr != nil {
// 			return nil, nil, inventoryListParseErr
// 		}
// 		inventoryDetails = append(inventoryDetails, inventoryDetailsTemp...)
// 	}

// 	return inventoryDetails, request.LastEvalutionKey, nil
// }

// func (repo *PurchasePersistance) GetTotalPurchaseOrders(request commonModels.InventoryListRequest) (int64, *commonModels.ErrorDetail) {
// 	var count int64 = 0
// 	proj := expression.NamesList(expression.Name("godownId"))
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey("")),
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
// 	result, getInventoryDetailError := getPurchaseDetails(expr, request)
// 	if getInventoryDetailError != nil {
// 		return 0, getInventoryDetailError
// 	}

// 	request.LastEvalutionKey = result.LastEvaluatedKey
// 	count = count + int64(len(result.Items))

// 	for len(result.Items) > 0 && request.LastEvalutionKey != nil {
// 		result, getInventoryDetailError = getPurchaseDetails(expr, request)
// 		if getInventoryDetailError != nil {
// 			return 0, getInventoryDetailError
// 		}
// 		request.LastEvalutionKey = result.LastEvaluatedKey
// 		count = count + int64(len(result.Items))
// 	}

// 	return count, nil
// }

// func (repo *PurchasePersistance) GetPurchaseBillDetails(request commonModels.InventoryFilterDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(request.GodownId)),
// 		expression.Key("inventorySortKey").BeginsWith(common.GetInventoryPurchanseSortKey(request.PurchaseBillNumber)),
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
// 	result, getInventoryDetailError := getPurchaseDetails(expr, commonModels.InventoryListRequest{
// 		InventoryFilterDto: request,
// 	})

// 	if getInventoryDetailError != nil {
// 		return nil, getInventoryDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		inventory, err := parseDbItemToInventory(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return inventory, nil
// 	}

// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("No purchase order found by the order no %s", request.PurchaseBillNumber),
// 	}
// }

// func (repo *PurchasePersistance) UpsertPurchaseOrder(data commonModels.InventoryDto) (*commonModels.InventoryDto, *commonModels.ErrorDetail) {
// 	av, err := dynamodbattribute.MarshalMap(data)
// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("Got error marshalling purchase details, purchase bill number - %s, godown id - %s, err: %s", data.BillNo, data.GodownId, err),
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
// 			ErrorMessage: fmt.Sprintf("Error in adding/updating purchase data for bill no %s, godown id %s, error message; %s", data.BillNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	return &data, nil
// }

// #region private functions
func getPurchaseDetails(expr expression.Expression, request commonModels.InventoryListRequest) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExclusiveStartKey:         request.LastEvalutionKey,
		TableName:                 aws.String(purchansePersistanceObj.purchaseTableName),
	}

	if request.PageSize > 0 {
		queryInput.Limit = &request.PageSize
	}
	result, err := purchansePersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get purchase bill for godown %s call failed: %s", request.GodownId, err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}

func getPurchaseDetailsInIndex(expr expression.Expression, request commonModels.InventoryListRequest, indexNane string) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExclusiveStartKey:         request.LastEvalutionKey,
		IndexName:                 &indexNane,
		TableName:                 aws.String(purchansePersistanceObj.purchaseTableName),
	}

	if request.PageSize > 0 {
		queryInput.Limit = &request.PageSize
	}
	result, err := purchansePersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get purchase bill for %+v call failed: %s", request, err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}
func parseDbItemsToPurchaseList(items []map[string]*dynamodb.AttributeValue) ([]commonModels.PurchaseDto, *commonModels.ErrorDetail) {
	purchaseDetails := make([]commonModels.PurchaseDto, 0)

	if len(items) > 0 {
		for _, val := range items {
			purchase, err := parseDbItemToPurchase(val)
			if err != nil {
				return nil, err
			}
			purchaseDetails = append(purchaseDetails, *purchase)
		}
		return purchaseDetails, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: "No data found",
	}
}
func parseDbItemToPurchase(item map[string]*dynamodb.AttributeValue) (*commonModels.ProductDto, *commonModels.ErrorDetail) {
	purchase := commonModels.ProductDto{}

	err := dynamodbattribute.UnmarshalMap(item, &purchase)

	if err != nil {
		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}

	return &purchase, nil
}

// #endregion

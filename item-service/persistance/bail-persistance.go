package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	jsii "github.com/aws/jsii-runtime-go"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var bailPersistanceObj *BailPersistance

type BailPersistance struct {
	db                 *dynamodb.DynamoDB
	inventoryTableName string
	bailInfoTable      string
	itemTable          string
}

func InitBailPersistance() (*BailPersistance, *commonModels.ErrorDetail) {
	if bailPersistanceObj == nil {
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

		bailPersistanceObj = &BailPersistance{
			db:                 dynamodb.New(dynamoDbSession),
			inventoryTableName: common.EnvValues.InventoryTableName,
			bailInfoTable:      common.EnvValues.BailInfoTableName,
			itemTable:          common.EnvValues.ItemTableName,
		}
	}

	return bailPersistanceObj, nil
}
func (repo *BailPersistance) UpdateBailDetailQuantity(data commonModels.BailDetailsDto) (*commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	updateItemInput := dynamodb.UpdateItemInput{
		TableName: &repo.itemTable,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(data.BranchId),
			},
			"sortKey": {
				S: aws.String(data.SortKey),
			},
		},
		UpdateExpression: jsii.String("SET receivedQuantity= :receivedQuantity, billedQuantity= :billedQuantity, pendingQuantity= :pendingQuantity"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":receivedQuantity": {
				N: jsii.String(fmt.Sprintf("%d", data.ReceivedQuantity)),
			},
			":billedQuantity": {
				N: jsii.String(fmt.Sprintf("%d", data.BilledQuantity)),
			},
			":pendingQuantity": {
				N: jsii.String(fmt.Sprintf("%d", data.PendingQuantity)),
			},
		},
	}
	_, err := repo.db.UpdateItem(&updateItemInput)

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in updating item quantities for bail no %s, branch id %s, error message; %s", data.BailNo, data.BranchId, err.Error()),
		}
	}
	return &data, nil

}
func (repo *BailPersistance) UpsertBailDetail(data commonModels.BailDetailsDto) (*commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling item details, bail number - %s, branch id - %s, err: %s", data.BailNo, data.BranchId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.itemTable,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating item data for bail no %s, branch id %s, error message; %s", data.BailNo, data.BranchId, err.Error()),
		}
	}
	return &data, nil
}

func (repo *BailPersistance) UpsertBailInfo(data commonModels.BailInfoDto) (*commonModels.BailInfoDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling Bail info, bail number - %s, branch id - %s, err: %s", data.BailNo, data.BranchId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.bailInfoTable,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating bail info for bail no %s, branch id %s, error message; %s", data.BailNo, data.BranchId, err.Error()),
		}
	}
	return &data, nil
}

func (repo *BailPersistance) GetBailInfoDetail(branchId, bailNo string) (*commonModels.BailInfoDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(branchId)),
		expression.Key("bailInfoSortKey").BeginsWith(common.GetBailInfoSortKey(bailNo)),
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
	result, getBailInfoDetailError := getBailInfoDetails(expr)

	if getBailInfoDetailError != nil {
		return nil, getBailInfoDetailError
	}

	if len(result.Items) > 0 {
		bailInfo, err := parseDbItemToBailInfo(result.Items[0])
		if err != nil {
			return nil, err
		}
		return bailInfo, nil
	}

	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: fmt.Sprintf("no bail found for bail number %s", bailNo),
	}
}
func (repo *BailPersistance) GetPurchasedBailDetailByQuanlity(branchId, quality string) ([]commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	var lastEvalutionKey map[string]*dynamodb.AttributeValue
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(branchId)),
		expression.Key("sortKey").BeginsWith(common.GetBailDetailPurchanseSortKey(quality, "")),
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
	result, getBailDetailError := getBailDetailsDetails(expr, lastEvalutionKey, 100)

	if getBailDetailError != nil {
		return nil, getBailDetailError
	}

	lastEvalutionKey = result.LastEvaluatedKey

	bailInfoDetails, parseListError := parseDbItemsToBailDetailList(result.Items)
	if parseListError != nil {
		return nil, parseListError
	}
	for len(result.Items) > 0 && result.LastEvaluatedKey != nil {
		result, getBailDetailError = getBailDetailsDetails(expr, lastEvalutionKey, 100)
		if getBailDetailError != nil {
			return nil, getBailDetailError
		}
		lastEvalutionKey = result.LastEvaluatedKey
		bailInfoDetailsTemp, inventoryListParseErr := parseDbItemsToBailDetailList(result.Items)
		if inventoryListParseErr != nil {
			return nil, inventoryListParseErr
		}
		bailInfoDetails = append(bailInfoDetails, bailInfoDetailsTemp...)

	}
	if len(bailInfoDetails) == 0 {
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: fmt.Sprintf("no bail found for quality %s", quality),
		}

	} else {
		return bailInfoDetails, nil
	}
}

func (repo *BailPersistance) GetOutofStockBailDetail(branchId, bailNo string) (*commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	bailInfo, gerBailInfoError := repo.GetBailInfoDetail(branchId, bailNo)
	if gerBailInfoError != nil {
		return nil, gerBailInfoError
	}

	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(branchId)),
		expression.Key("sortKey").BeginsWith(common.GetBailDetailOutOfStockSortKey(bailInfo.Quality, bailInfo.BailNo)),
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
	result, getBailDetailError := getBailDetailsDetails(expr, nil, 0)

	if getBailDetailError != nil {
		return nil, getBailDetailError
	}

	if len(result.Items) > 0 {
		bailInfo, err := parseDbItemToBailDetails(result.Items[0])
		if err != nil {
			return nil, err
		}
		return bailInfo, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: fmt.Sprintf("no bail found for bail number %s", bailNo),
	}
}

func (repo *BailPersistance) GetPurchasedBailDetail(branchId, bailNo string) (*commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	bailInfo, gerBailInfoError := repo.GetBailInfoDetail(branchId, bailNo)
	if gerBailInfoError != nil {
		return nil, gerBailInfoError
	}

	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(branchId)),
		expression.Key("sortKey").BeginsWith(common.GetBailDetailPurchanseSortKey(bailInfo.Quality, bailInfo.BailNo)),
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
	result, getBailDetailError := getBailDetailsDetails(expr, nil, 0)

	if getBailDetailError != nil {
		return nil, getBailDetailError
	}

	if len(result.Items) > 0 {
		bailInfo, err := parseDbItemToBailDetails(result.Items[0])
		if err != nil {
			return nil, err
		}
		return bailInfo, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: fmt.Sprintf("no bail found for bail number %s", bailNo),
	}
}

func (repo *BailPersistance) GetSalesBailDetail(branchId, bailNo, salesBillNumber string) ([]commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	var lastEvalutionKey map[string]*dynamodb.AttributeValue
	var pageSize int64 = 10
	bailInfo, gerBailInfoError := repo.GetBailInfoDetail(branchId, bailNo)
	if gerBailInfoError != nil {
		return nil, gerBailInfoError
	}
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(branchId)),
		expression.Key("sortKey").BeginsWith(common.GetBailDetailSalesSortKey(bailInfo.Quality, bailInfo.BailNo, salesBillNumber)),
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
	result, getbailDetailError := getBailDetailsDetails(expr, nil, 0)
	if getbailDetailError != nil {
		return nil, getbailDetailError
	}

	lastEvalutionKey = result.LastEvaluatedKey
	bailDetails, inventoryListParseErr := parseDbItemsToBailDetailList(result.Items)
	if inventoryListParseErr != nil {
		return nil, inventoryListParseErr
	}
	for len(bailDetails) < int(pageSize) && lastEvalutionKey != nil {
		result, getbailDetailError = getBailDetailsDetails(expr, nil, 0)
		if getbailDetailError != nil {
			return nil, getbailDetailError
		}
		lastEvalutionKey = result.LastEvaluatedKey
		bailDetailsTemp, bailListParseErr := parseDbItemsToBailDetailList(result.Items)
		if bailListParseErr != nil {
			return nil, bailListParseErr
		}
		bailDetails = append(bailDetails, bailDetailsTemp...)
	}
	return bailDetails, nil
}

func (repo *BailPersistance) DeleteBailInfo(branchId, bailno string) *commonModels.ErrorDetail {
	bailInfo, getBailInfoErr := repo.GetBailInfoDetail(branchId, bailno)
	if getBailInfoErr != nil {
		return getBailInfoErr
	}
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.bailInfoTable,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"bailInfoSortKey": {
				S: aws.String(bailInfo.BailInfoSortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting bail no: %s for branchId: %s, error; %s", bailno, branchId, err.Error()),
		}
	}
	timeNow := time.Now().UTC().Unix()
	bailInfo.BailInfoSortKey = fmt.Sprintf("%s|%d", common.GetBailInfoDeleteSortKey(bailno), timeNow)
	_, updateError := repo.UpsertBailInfo(*bailInfo)
	if updateError != nil {
		return updateError
	}

	return nil
}

func (repo *BailPersistance) DeleteSalesBailDetails(branchId, bailno, salesBillNumber string) *commonModels.ErrorDetail {
	itemInfo, getBailInfoErr := repo.GetSalesBailDetail(branchId, bailno, salesBillNumber)
	if getBailInfoErr != nil {
		return getBailInfoErr
	}
	for _, val := range itemInfo {
		_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: &repo.itemTable,
			Key: map[string]*dynamodb.AttributeValue{
				"branchId": {
					S: aws.String(branchId),
				},
				"sortKey": {
					S: aws.String(val.SortKey),
				},
			},
		})

		if err != nil {
			common.WriteLog(1, err.Error())
			return &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorDelete,
				ErrorMessage: fmt.Sprintf("Error in deleting bail no: %s for branchId: %s, error; %s", bailno, branchId, err.Error()),
			}
		}
		timeNow := time.Now().UTC().Unix()
		val.SortKey = fmt.Sprintf("%s|%d", common.GetBailDetailDeleteSortKey(val.Quality, val.BailNo), timeNow)
		_, updateError := repo.UpsertBailDetail(val)
		if updateError != nil {
			return updateError
		}
	}
	return nil
}

func (repo *BailPersistance) DeleteBailDetails(branchId, bailno string) *commonModels.ErrorDetail {
	itemInfo, getBailInfoErr := repo.GetPurchasedBailDetail(branchId, bailno)
	if getBailInfoErr != nil {
		return getBailInfoErr
	}
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.itemTable,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(itemInfo.SortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting bail no: %s for branchId: %s, error; %s", bailno, branchId, err.Error()),
		}
	}
	timeNow := time.Now().UTC().Unix()
	itemInfo.SortKey = fmt.Sprintf("%s|%d", common.GetBailDetailDeleteSortKey(itemInfo.Quality, itemInfo.BailNo), timeNow)
	_, updateError := repo.UpsertBailDetail(*itemInfo)
	if updateError != nil {
		return updateError
	}

	return nil
}
func (repo *BailPersistance) RegenrateOutofStockBail(data commonModels.BailDetailsDto) *commonModels.ErrorDetail {
	itemInfo, getBailInfoErr := repo.GetOutofStockBailDetail(data.BranchId, data.BailNo)
	if getBailInfoErr != nil {
		return getBailInfoErr
	}
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.itemTable,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(data.BranchId),
			},
			"sortKey": {
				S: aws.String(itemInfo.SortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting bail no: %s for branchId: %s, error; %s", data.BailNo, data.BranchId, err.Error()),
		}
	}
	itemInfo.SortKey = common.GetBailDetailPurchanseSortKey(itemInfo.Quality, itemInfo.BailNo)
	itemInfo.PendingQuantity = data.PendingQuantity
	_, updateError := repo.UpsertBailDetail(*itemInfo)
	if updateError != nil {
		return updateError
	}

	return nil
}
func (repo *BailPersistance) UpdateBailDetailsOutofStock(branchId, bailno string) *commonModels.ErrorDetail {
	itemInfo, getBailInfoErr := repo.GetPurchasedBailDetail(branchId, bailno)
	if getBailInfoErr != nil {
		return getBailInfoErr
	}
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.itemTable,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(itemInfo.SortKey),
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting bail no: %s for branchId: %s, error; %s", bailno, branchId, err.Error()),
		}
	}
	itemInfo.SortKey = common.GetBailDetailOutOfStockSortKey(itemInfo.Quality, itemInfo.BailNo)
	itemInfo.PendingQuantity = 0
	_, updateError := repo.UpsertBailDetail(*itemInfo)
	if updateError != nil {
		return updateError
	}

	return nil
}

// #region private functions
func getBailInfoDetails(expr expression.Expression) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	result, err := purchansePersistanceObj.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(purchansePersistanceObj.bailInfoTable),
	})

	if err != nil {
		errMessage := fmt.Sprintf("get bail info for, call failed err: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}

func parseDbItemToBailInfo(item map[string]*dynamodb.AttributeValue) (*commonModels.BailInfoDto, *commonModels.ErrorDetail) {
	bailInfo := commonModels.BailInfoDto{}

	err := dynamodbattribute.UnmarshalMap(item, &bailInfo)

	if err != nil {
		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}

	return &bailInfo, nil
}

func getBailDetailsDetails(expr expression.Expression, lastEvalutionKey map[string]*dynamodb.AttributeValue, pageSize int64) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
	var queryInput = dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExclusiveStartKey:         lastEvalutionKey,
		TableName:                 aws.String(purchansePersistanceObj.itemTable),
	}

	if pageSize > 0 {
		queryInput.Limit = &pageSize
	}

	result, err := purchansePersistanceObj.db.Query(&queryInput)

	if err != nil {
		errMessage := fmt.Sprintf("get bail info for, call failed err: %s", err.Error())
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}
	return result, nil
}

func parseDbItemsToBailDetailList(items []map[string]*dynamodb.AttributeValue) ([]commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	bailDetails := make([]commonModels.BailDetailsDto, 0)

	if len(items) > 0 {
		for _, val := range items {
			bailDetail, err := parseDbItemToBailDetails(val)
			if err != nil {
				return nil, err
			}
			bailDetails = append(bailDetails, *bailDetail)
		}
		return bailDetails, nil
	}
	return nil, &commonModels.ErrorDetail{
		ErrorCode:    commonModels.ErrorNoDataFound,
		ErrorMessage: "No data found",
	}
}

func parseDbItemToBailDetails(item map[string]*dynamodb.AttributeValue) (*commonModels.BailDetailsDto, *commonModels.ErrorDetail) {
	bailDetail := commonModels.BailDetailsDto{}

	err := dynamodbattribute.UnmarshalMap(item, &bailDetail)

	if err != nil {
		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
		common.WriteLog(1, errMessage)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: errMessage,
		}
	}

	return &bailDetail, nil
}

// #endregion

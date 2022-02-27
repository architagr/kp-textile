package persistance

// var balePersistanceObj *BalePersistance

// type BalePersistance struct {
// 	db                 *dynamodb.DynamoDB
// 	inventoryTableName string
// 	baleInfoTable      string
// 	itemTable          string
// }

// func InitBalePersistance() (*BalePersistance, *commonModels.ErrorDetail) {
// 	if balePersistanceObj == nil {
// 		dbSession, err := session.NewSessionWithOptions(session.Options{
// 			SharedConfigState: session.SharedConfigEnable,
// 		})

// 		if err != nil {
// 			return nil, &commonModels.ErrorDetail{
// 				ErrorCode:    commonModels.ErrorDbConnection,
// 				ErrorMessage: err.Error(),
// 			}
// 		}
// 		dynamoDbSession := session.Must(dbSession, err)

// 		balePersistanceObj = &BalePersistance{
// 			db:                 dynamodb.New(dynamoDbSession),
// 			inventoryTableName: common.EnvValues.InventoryTableName,
// 			baleInfoTable:      common.EnvValues.BaleInfoTableName,
// 			itemTable:          common.EnvValues.ItemTableName,
// 		}
// 	}

// 	return balePersistanceObj, nil
// }
// func (repo *BalePersistance) UpdateBaleDetailQuantity(data commonModels.BaleDetailsDto) (*commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	updateItemInput := dynamodb.UpdateItemInput{
// 		TableName: &repo.itemTable,
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"godownId": {
// 				S: aws.String(data.GodownId),
// 			},
// 			"sortKey": {
// 				S: aws.String(data.SortKey),
// 			},
// 		},
// 		UpdateExpression: jsii.String("SET receivedQuantity= :receivedQuantity, billedQuantity= :billedQuantity, pendingQuantity= :pendingQuantity"),
// 		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
// 			":receivedQuantity": {
// 				N: jsii.String(fmt.Sprintf("%d", data.ReceivedQuantity)),
// 			},
// 			":billedQuantity": {
// 				N: jsii.String(fmt.Sprintf("%d", data.BilledQuantity)),
// 			},
// 			":pendingQuantity": {
// 				N: jsii.String(fmt.Sprintf("%d", data.PendingQuantity)),
// 			},
// 		},
// 	}
// 	_, err := repo.db.UpdateItem(&updateItemInput)

// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorInsert,
// 			ErrorMessage: fmt.Sprintf("Error in updating item quantities for bale no %s, branch id %s, error message; %s", data.BaleNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	return &data, nil

// }
// func (repo *BalePersistance) UpsertBaleDetail(data commonModels.BaleDetailsDto) (*commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	av, err := dynamodbattribute.MarshalMap(data)
// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("Got error marshalling item details, bale number - %s, branch id - %s, err: %s", data.BaleNo, data.GodownId, err),
// 		}
// 	}
// 	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
// 		TableName: &repo.itemTable,
// 		Item:      av,
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorInsert,
// 			ErrorMessage: fmt.Sprintf("Error in adding/updating item data for bale no %s, branch id %s, error message; %s", data.BaleNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	return &data, nil
// }

// func (repo *BalePersistance) UpsertBaleInfo(data commonModels.BaleInfoDto) (*commonModels.BaleInfoDto, *commonModels.ErrorDetail) {
// 	av, err := dynamodbattribute.MarshalMap(data)
// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("Got error marshalling Bale info, bale number - %s, branch id - %s, err: %s", data.BaleNo, data.GodownId, err),
// 		}
// 	}
// 	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
// 		TableName: &repo.baleInfoTable,
// 		Item:      av,
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())

// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorInsert,
// 			ErrorMessage: fmt.Sprintf("Error in adding/updating bale info for bale no %s, branch id %s, error message; %s", data.BaleNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	return &data, nil
// }

// func (repo *BalePersistance) GetBaleInfoDetail(godownId, baleNo string) (*commonModels.BaleInfoDto, *commonModels.ErrorDetail) {
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(godownId)),
// 		expression.Key("baleInfoSortKey").BeginsWith(common.GetBaleInfoSortKey(baleNo)),
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
// 	result, getBaleInfoDetailError := getBaleInfoDetails(expr)

// 	if getBaleInfoDetailError != nil {
// 		return nil, getBaleInfoDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		baleInfo, err := parseDbItemToBaleInfo(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return baleInfo, nil
// 	}

// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("no bale found for bale number %s", baleNo),
// 	}
// }
// func (repo *BalePersistance) GetPurchasedBaleDetailByQuanlity(godownId, quality string) ([]commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	var lastEvalutionKey map[string]*dynamodb.AttributeValue
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(godownId)),
// 		expression.Key("sortKey").BeginsWith(common.GetBaleDetailPurchanseSortKey(quality, "")),
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
// 	result, getBaleDetailError := getBaleDetailsDetails(expr, lastEvalutionKey, 100)

// 	if getBaleDetailError != nil {
// 		return nil, getBaleDetailError
// 	}

// 	lastEvalutionKey = result.LastEvaluatedKey

// 	baleInfoDetails, parseListError := parseDbItemsToBaleDetailList(result.Items)
// 	if parseListError != nil {
// 		return nil, parseListError
// 	}
// 	for len(result.Items) > 0 && result.LastEvaluatedKey != nil {
// 		result, getBaleDetailError = getBaleDetailsDetails(expr, lastEvalutionKey, 100)
// 		if getBaleDetailError != nil {
// 			return nil, getBaleDetailError
// 		}
// 		lastEvalutionKey = result.LastEvaluatedKey
// 		baleInfoDetailsTemp, inventoryListParseErr := parseDbItemsToBaleDetailList(result.Items)
// 		if inventoryListParseErr != nil {
// 			return nil, inventoryListParseErr
// 		}
// 		baleInfoDetails = append(baleInfoDetails, baleInfoDetailsTemp...)

// 	}
// 	if len(baleInfoDetails) == 0 {
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorNoDataFound,
// 			ErrorMessage: fmt.Sprintf("no bale found for quality %s", quality),
// 		}

// 	} else {
// 		return baleInfoDetails, nil
// 	}
// }

// func (repo *BalePersistance) GetOutofStockBaleDetail(godownId, baleNo string) (*commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	baleInfo, gerBaleInfoError := repo.GetBaleInfoDetail(godownId, baleNo)
// 	if gerBaleInfoError != nil {
// 		return nil, gerBaleInfoError
// 	}

// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(godownId)),
// 		expression.Key("sortKey").BeginsWith(common.GetBaleDetailOutOfStockSortKey(baleInfo.Quality, baleInfo.BaleNo)),
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
// 	result, getBaleDetailError := getBaleDetailsDetails(expr, nil, 0)

// 	if getBaleDetailError != nil {
// 		return nil, getBaleDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		baleInfo, err := parseDbItemToBaleDetails(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return baleInfo, nil
// 	}
// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("no bale found for bale number %s", baleNo),
// 	}
// }

// func (repo *BalePersistance) GetPurchasedBaleDetail(godownId, baleNo string) (*commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	baleInfo, gerBaleInfoError := repo.GetBaleInfoDetail(godownId, baleNo)
// 	if gerBaleInfoError != nil {
// 		return nil, gerBaleInfoError
// 	}

// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(godownId)),
// 		expression.Key("sortKey").BeginsWith(common.GetBaleDetailPurchanseSortKey(baleInfo.Quality, baleInfo.BaleNo)),
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
// 	result, getBaleDetailError := getBaleDetailsDetails(expr, nil, 0)

// 	if getBaleDetailError != nil {
// 		return nil, getBaleDetailError
// 	}

// 	if len(result.Items) > 0 {
// 		baleInfo, err := parseDbItemToBaleDetails(result.Items[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		return baleInfo, nil
// 	}
// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: fmt.Sprintf("no bale found for bale number %s", baleNo),
// 	}
// }

// func (repo *BalePersistance) GetSalesBaleDetail(godownId, baleNo, salesBillNumber string) ([]commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	var lastEvalutionKey map[string]*dynamodb.AttributeValue
// 	var pageSize int64 = 10
// 	baleInfo, gerBaleInfoError := repo.GetBaleInfoDetail(godownId, baleNo)
// 	if gerBaleInfoError != nil {
// 		return nil, gerBaleInfoError
// 	}
// 	keyCondition := expression.KeyAnd(
// 		expression.Key("godownId").Equal(expression.Value(godownId)),
// 		expression.Key("sortKey").BeginsWith(common.GetBaleDetailSalesSortKey(baleInfo.Quality, baleInfo.BaleNo, salesBillNumber)),
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
// 	result, getbaleDetailError := getBaleDetailsDetails(expr, nil, 0)
// 	if getbaleDetailError != nil {
// 		return nil, getbaleDetailError
// 	}

// 	lastEvalutionKey = result.LastEvaluatedKey
// 	baleDetails, inventoryListParseErr := parseDbItemsToBaleDetailList(result.Items)
// 	if inventoryListParseErr != nil {
// 		return nil, inventoryListParseErr
// 	}
// 	for len(baleDetails) < int(pageSize) && lastEvalutionKey != nil {
// 		result, getbaleDetailError = getBaleDetailsDetails(expr, nil, 0)
// 		if getbaleDetailError != nil {
// 			return nil, getbaleDetailError
// 		}
// 		lastEvalutionKey = result.LastEvaluatedKey
// 		baleDetailsTemp, baleListParseErr := parseDbItemsToBaleDetailList(result.Items)
// 		if baleListParseErr != nil {
// 			return nil, baleListParseErr
// 		}
// 		baleDetails = append(baleDetails, baleDetailsTemp...)
// 	}
// 	return baleDetails, nil
// }

// func (repo *BalePersistance) DeleteBaleInfo(godownId, baleno string) *commonModels.ErrorDetail {
// 	baleInfo, getBaleInfoErr := repo.GetBaleInfoDetail(godownId, baleno)
// 	if getBaleInfoErr != nil {
// 		return getBaleInfoErr
// 	}
// 	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
// 		TableName: &repo.baleInfoTable,
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"godownId": {
// 				S: aws.String(godownId),
// 			},
// 			"baleInfoSortKey": {
// 				S: aws.String(baleInfo.BaleInfoSortKey),
// 			},
// 		},
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorDelete,
// 			ErrorMessage: fmt.Sprintf("Error in deleting bale no: %s for godownId: %s, error; %s", baleno, godownId, err.Error()),
// 		}
// 	}
// 	timeNow := time.Now().UTC().Unix()
// 	baleInfo.BaleInfoSortKey = fmt.Sprintf("%s|%d", common.GetBaleInfoDeleteSortKey(baleno), timeNow)
// 	_, updateError := repo.UpsertBaleInfo(*baleInfo)
// 	if updateError != nil {
// 		return updateError
// 	}

// 	return nil
// }

// func (repo *BalePersistance) DeleteSalesBaleDetails(godownId, baleno, salesBillNumber string) *commonModels.ErrorDetail {
// 	itemInfo, getBaleInfoErr := repo.GetSalesBaleDetail(godownId, baleno, salesBillNumber)
// 	if getBaleInfoErr != nil {
// 		return getBaleInfoErr
// 	}
// 	for _, val := range itemInfo {
// 		_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
// 			TableName: &repo.itemTable,
// 			Key: map[string]*dynamodb.AttributeValue{
// 				"godownId": {
// 					S: aws.String(godownId),
// 				},
// 				"sortKey": {
// 					S: aws.String(val.SortKey),
// 				},
// 			},
// 		})

// 		if err != nil {
// 			common.WriteLog(1, err.Error())
// 			return &commonModels.ErrorDetail{
// 				ErrorCode:    commonModels.ErrorDelete,
// 				ErrorMessage: fmt.Sprintf("Error in deleting bale no: %s for godownId: %s, error; %s", baleno, godownId, err.Error()),
// 			}
// 		}
// 		timeNow := time.Now().UTC().Unix()
// 		val.SortKey = fmt.Sprintf("%s|%d", common.GetBaleDetailDeleteSortKey(val.Quality, val.BaleNo), timeNow)
// 		_, updateError := repo.UpsertBaleDetail(val)
// 		if updateError != nil {
// 			return updateError
// 		}
// 	}
// 	return nil
// }

// func (repo *BalePersistance) DeleteBaleDetails(godownId, baleno string) *commonModels.ErrorDetail {
// 	itemInfo, getBaleInfoErr := repo.GetPurchasedBaleDetail(godownId, baleno)
// 	if getBaleInfoErr != nil {
// 		return getBaleInfoErr
// 	}
// 	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
// 		TableName: &repo.itemTable,
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"godownId": {
// 				S: aws.String(godownId),
// 			},
// 			"sortKey": {
// 				S: aws.String(itemInfo.SortKey),
// 			},
// 		},
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorDelete,
// 			ErrorMessage: fmt.Sprintf("Error in deleting bale no: %s for godownId: %s, error; %s", baleno, godownId, err.Error()),
// 		}
// 	}
// 	timeNow := time.Now().UTC().Unix()
// 	itemInfo.SortKey = fmt.Sprintf("%s|%d", common.GetBaleDetailDeleteSortKey(itemInfo.Quality, itemInfo.BaleNo), timeNow)
// 	_, updateError := repo.UpsertBaleDetail(*itemInfo)
// 	if updateError != nil {
// 		return updateError
// 	}

// 	return nil
// }
// func (repo *BalePersistance) RegenrateOutofStockBale(data commonModels.BaleDetailsDto) *commonModels.ErrorDetail {
// 	itemInfo, getBaleInfoErr := repo.GetOutofStockBaleDetail(data.GodownId, data.BaleNo)
// 	if getBaleInfoErr != nil {
// 		return getBaleInfoErr
// 	}
// 	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
// 		TableName: &repo.itemTable,
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"godownId": {
// 				S: aws.String(data.GodownId),
// 			},
// 			"sortKey": {
// 				S: aws.String(itemInfo.SortKey),
// 			},
// 		},
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorDelete,
// 			ErrorMessage: fmt.Sprintf("Error in deleting bale no: %s for godownId: %s, error; %s", data.BaleNo, data.GodownId, err.Error()),
// 		}
// 	}
// 	itemInfo.SortKey = common.GetBaleDetailPurchanseSortKey(itemInfo.Quality, itemInfo.BaleNo)
// 	itemInfo.PendingQuantity = data.PendingQuantity
// 	_, updateError := repo.UpsertBaleDetail(*itemInfo)
// 	if updateError != nil {
// 		return updateError
// 	}

// 	return nil
// }
// func (repo *BalePersistance) UpdateBaleDetailsOutofStock(godownId, baleno string) *commonModels.ErrorDetail {
// 	itemInfo, getBaleInfoErr := repo.GetPurchasedBaleDetail(godownId, baleno)
// 	if getBaleInfoErr != nil {
// 		return getBaleInfoErr
// 	}
// 	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
// 		TableName: &repo.itemTable,
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"godownId": {
// 				S: aws.String(godownId),
// 			},
// 			"sortKey": {
// 				S: aws.String(itemInfo.SortKey),
// 			},
// 		},
// 	})

// 	if err != nil {
// 		common.WriteLog(1, err.Error())
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorDelete,
// 			ErrorMessage: fmt.Sprintf("Error in deleting bale no: %s for godownId: %s, error; %s", baleno, godownId, err.Error()),
// 		}
// 	}
// 	itemInfo.SortKey = common.GetBaleDetailOutOfStockSortKey(itemInfo.Quality, itemInfo.BaleNo)
// 	itemInfo.PendingQuantity = 0
// 	_, updateError := repo.UpsertBaleDetail(*itemInfo)
// 	if updateError != nil {
// 		return updateError
// 	}

// 	return nil
// }

// // #region private functions
// func getBaleInfoDetails(expr expression.Expression) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
// 	result, err := purchansePersistanceObj.db.Query(&dynamodb.QueryInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		FilterExpression:          expr.Filter(),
// 		KeyConditionExpression:    expr.KeyCondition(),
// 		TableName:                 aws.String(purchansePersistanceObj.baleInfoTable),
// 	})

// 	if err != nil {
// 		errMessage := fmt.Sprintf("get bale info for, call failed err: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	return result, nil
// }

// func parseDbItemToBaleInfo(item map[string]*dynamodb.AttributeValue) (*commonModels.BaleInfoDto, *commonModels.ErrorDetail) {
// 	baleInfo := commonModels.BaleInfoDto{}

// 	err := dynamodbattribute.UnmarshalMap(item, &baleInfo)

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}

// 	return &baleInfo, nil
// }

// func getBaleDetailsDetails(expr expression.Expression, lastEvalutionKey map[string]*dynamodb.AttributeValue, pageSize int64) (*dynamodb.QueryOutput, *commonModels.ErrorDetail) {
// 	var queryInput = dynamodb.QueryInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		FilterExpression:          expr.Filter(),
// 		KeyConditionExpression:    expr.KeyCondition(),
// 		ExclusiveStartKey:         lastEvalutionKey,
// 		TableName:                 aws.String(purchansePersistanceObj.itemTable),
// 	}

// 	if pageSize > 0 {
// 		queryInput.Limit = &pageSize
// 	}

// 	result, err := purchansePersistanceObj.db.Query(&queryInput)

// 	if err != nil {
// 		errMessage := fmt.Sprintf("get bale info for, call failed err: %s", err.Error())
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}
// 	return result, nil
// }

// func parseDbItemsToBaleDetailList(items []map[string]*dynamodb.AttributeValue) ([]commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	baleDetails := make([]commonModels.BaleDetailsDto, 0)

// 	if len(items) > 0 {
// 		for _, val := range items {
// 			baleDetail, err := parseDbItemToBaleDetails(val)
// 			if err != nil {
// 				return nil, err
// 			}
// 			baleDetails = append(baleDetails, *baleDetail)
// 		}
// 		return baleDetails, nil
// 	}
// 	return nil, &commonModels.ErrorDetail{
// 		ErrorCode:    commonModels.ErrorNoDataFound,
// 		ErrorMessage: "No data found",
// 	}
// }

// func parseDbItemToBaleDetails(item map[string]*dynamodb.AttributeValue) (*commonModels.BaleDetailsDto, *commonModels.ErrorDetail) {
// 	baleDetail := commonModels.BaleDetailsDto{}

// 	err := dynamodbattribute.UnmarshalMap(item, &baleDetail)

// 	if err != nil {
// 		errMessage := fmt.Sprintf("Got error unmarshalling: %s", err)
// 		common.WriteLog(1, errMessage)
// 		return nil, &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: errMessage,
// 		}
// 	}

// 	return &baleDetail, nil
// }

// // #endregion

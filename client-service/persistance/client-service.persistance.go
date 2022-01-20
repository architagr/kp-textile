package persistance

import (
	"client-service/common"
	commonModels "commonpkg/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var clientServicePersistanceObj *ClientServicePersistance

type ClientServicePersistance struct {
	db              *dynamodb.DynamoDB
	clientTableName string
}

func InitClientServicePersistance() (*ClientServicePersistance, *commonModels.ErrorDetail) {
	if clientServicePersistanceObj == nil {
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

		clientServicePersistanceObj = &ClientServicePersistance{
			db:              dynamodb.New(dynamoDbSession),
			clientTableName: common.EnvValues.ClientTableName,
		}
	}

	return clientServicePersistanceObj, nil
}

func (repo *ClientServicePersistance) GetPersonByClientId(request commonModels.GetClientRequestDto) ([]commonModels.ContactPersonDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").BeginsWith(common.GetClientContactSortKey(request.ClientId, "")),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("Got error building expression: %s", err.Error()))
	}

	result, err := repo.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.clientTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get person for client %s call failed: %s", request.ClientId, err.Error()))
	}

	clientPersons := make([]commonModels.ContactPersonDto, len(result.Items))

	for i, val := range result.Items {
		clientPerson := commonModels.ContactPersonDto{}

		err = dynamodbattribute.UnmarshalMap(val, &clientPerson)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		clientPersons[i] = clientPerson
	}
	return clientPersons, nil
}

func (repo *ClientServicePersistance) GetClient(request commonModels.GetClientRequestDto) (commonModels.ClientDto, *commonModels.ErrorDetail) {

	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").Equal(expression.Value(common.GetClientSortKey(request.ClientId))),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("Got error building expression: %s", err.Error()))
	}

	result, err := repo.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.clientTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get client call failed: %s", err.Error()))
	}

	client := commonModels.ClientDto{}
	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &client)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
	}
	return client, nil
}

func buildFilterExpression(filterData commonModels.ClientListRequest, projection *expression.ProjectionBuilder) (*expression.Expression, *commonModels.ErrorDetail) {

	filter := expression.Name("branchId").Equal(expression.Value(filterData.BranchId)).And(expression.Name("sortKey").BeginsWith(common.ClientSortKey))

	if len(filterData.Alias) > 0 {
		filter = filter.And(expression.Name("alias").Contains(filterData.Alias))
	}

	if len(filterData.CompanyName) > 0 {
		filter = filter.And(expression.Name("companyName").Contains(filterData.CompanyName))
	}

	if len(filterData.Email) > 0 {
		filter = filter.And(expression.Name("contactInfo.email").Contains(filterData.Email))
	}

	if len(filterData.ContactPersonFirstName) > 0 {
		filter = filter.And(expression.Name("contactPersons.firstName").Contains(filterData.ContactPersonFirstName))
	}

	if len(filterData.ContactPersonFirstName) > 0 {
		filter = filter.And(expression.Name("contactPersons.lastName").Contains(filterData.ContactPersonFirstName))
	}

	if len(filterData.PaymentTerm) > 0 {
		filter = filter.And(expression.Name("paymentTerm").Contains(filterData.PaymentTerm))
	}
	builder := expression.NewBuilder().WithFilter(filter)

	if projection != nil {
		builder.WithProjection(*projection)
	}
	expr, err := builder.Build()

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("Got error building expression: %s", err.Error()))
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInvalidRequestParam,
			ErrorMessage: "Error building filter",
		}
	}

	return &expr, nil
}

func (repo *ClientServicePersistance) GetClientTotalByFilter(filterData commonModels.ClientListRequest) (int64, *commonModels.ErrorDetail) {
	var count int64
	proj := expression.NamesList(expression.Name("branchId"))
	expr, errorDetails := buildFilterExpression(filterData, &proj)
	if errorDetails != nil {
		return count, errorDetails
	}

	var result *dynamodb.ScanOutput
	var err error

	result, err = repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.clientTableName),
		ExclusiveStartKey:         filterData.LastEvalutionKey,
		Limit:                     aws.Int64(filterData.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter client call failed: %s", err.Error()))
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))
	for len(result.Items) > 0 && result.LastEvaluatedKey != nil {
		result, err = repo.db.Scan(&dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 aws.String(repo.clientTableName),
			ExclusiveStartKey:         filterData.LastEvalutionKey,
			Limit:                     aws.Int64(filterData.PageSize),
		})
		if err != nil {
			common.WriteLog(1, fmt.Sprintf("filter client call failed: %s", err.Error()))
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}
	return count, nil
}
func (repo *ClientServicePersistance) scanClient(expr *expression.Expression, filterdata *commonModels.ClientListRequest) (*dynamodb.ScanOutput, *commonModels.ErrorDetail) {
	result, err := repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.clientTableName),
		ExclusiveStartKey:         filterdata.LastEvalutionKey,
		Limit:                     aws.Int64(filterdata.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter client call failed: %s", err.Error()))
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: err.Error(),
		}
	}
	return result, nil
}

func parseClientScanResult(reultItem []map[string]*dynamodb.AttributeValue) []commonModels.ClientDto {
	clients := make([]commonModels.ClientDto, 0)
	if len(reultItem) > 0 {
		for _, val := range reultItem {
			client := commonModels.ClientDto{}

			err := dynamodbattribute.UnmarshalMap(val, &client)

			if err != nil {
				log.Fatalf("Got error unmarshalling: %s", err)
			}
			clients = append(clients, client)
		}
	}
	return clients
}
func (repo *ClientServicePersistance) GetClientByFilter(filterData commonModels.ClientListRequest) ([]commonModels.ClientDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {

	expr, errorDetails := buildFilterExpression(filterData, nil)
	if errorDetails != nil {
		return nil, nil, errorDetails
	}
	result, errorDetails := repo.scanClient(expr, &filterData)

	if errorDetails != nil {
		return nil, nil, errorDetails
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	clients := parseClientScanResult(result.Items)

	for len(clients) < int(filterData.PageSize) && filterData.LastEvalutionKey != nil {

		result, errorDetails = repo.scanClient(expr, &filterData)
		if errorDetails != nil {
			return nil, nil, errorDetails
		}

		clientsTemp := parseClientScanResult(result.Items)
		if len(clientsTemp) > 0 {
			clients = append(clients, clientsTemp...)
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
	}

	return clients, filterData.LastEvalutionKey, nil
}

func (repo *ClientServicePersistance) UpsertClient(client commonModels.ClientDto, isNew bool) (*commonModels.ClientDto, *commonModels.ErrorDetail) {
	existigClients, _, errorDetails := repo.GetClientByFilter(commonModels.ClientListRequest{
		ClientFilterDto: commonModels.ClientFilterDto{BranchId: client.BranchId,
			CompanyName: client.CompanyName,
			Alias:       client.Alias,
			Email:       client.ContactInfo.Email,
		},
		PageSize: 10,
	})

	if errorDetails != nil {
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Error in validating exiting client, error: %s", errorDetails.Error()),
		}
	}
	if len(existigClients) > 0 {
		flag := true
		if !isNew {
			flag = false
			for _, val := range existigClients {
				if val.ClientId != client.ClientId {
					flag = true
				}
			}
		}
		if flag {
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorInsert,
				ErrorMessage: "similar client already exists",
			}
		}
	}
	av, err := dynamodbattribute.MarshalMap(client)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling client details item, client name - %s, client id - %s, err: %s", client.CompanyName, client.ClientId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.clientTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating Client %s, client id %s, error message; %s", client.CompanyName, client.ClientId, err.Error()),
		}
	}
	return &client, nil
}

func (repo *ClientServicePersistance) UpsertClientContact(clientContact commonModels.ContactPersonDto) (*commonModels.ContactPersonDto, *commonModels.ErrorDetail) {

	av, err := dynamodbattribute.MarshalMap(clientContact)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new Client Contact detailes clinet id: %s, err: %s", clientContact.ClientId, err.Error()),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.clientTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding Client contact (%s) for client id %s, error message; %s", clientContact.FirstName, clientContact.ClientId, err.Error()),
		}
	}
	return &clientContact, nil
}

func (repo *ClientServicePersistance) DeleteClientContact(branchId, clientId, contactId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.clientTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetClientContactSortKey(clientId, contactId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())

		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting Client contact id (%s) for client id %s, error message; %s", contactId, clientId, err.Error()),
		}
	}
	return nil
}

func (repo *ClientServicePersistance) DeleteClient(branchId, clientId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.clientTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetClientSortKey(clientId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting client id %s, error message; %s", clientId, err.Error()),
		}
	}
	return nil
}

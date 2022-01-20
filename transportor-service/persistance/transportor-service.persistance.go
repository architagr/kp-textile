package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"log"
	"transportor-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var transporterPersistanceObj *TransporterPersistance

type TransporterPersistance struct {
	db                   *dynamodb.DynamoDB
	transporterTableName string
}

func InitTransporterPersistance() (*TransporterPersistance, *commonModels.ErrorDetail) {
	if transporterPersistanceObj == nil {
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

		transporterPersistanceObj = &TransporterPersistance{
			db:                   dynamodb.New(dynamoDbSession),
			transporterTableName: common.EnvValues.TransporterTableName,
		}
	}

	return transporterPersistanceObj, nil
}

func (repo *TransporterPersistance) GetPersonByTransporterId(request commonModels.GetTransporterRequestDto) ([]commonModels.TransporterContactPersonDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").BeginsWith(common.GetTransporterContactSortKey(request.TransporterId, "")),
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
		TableName:                 aws.String(repo.transporterTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get person for transporter %s call failed: %s", request.TransporterId, err.Error()))
	}

	transporterPersons := make([]commonModels.TransporterContactPersonDto, len(result.Items))

	for i, val := range result.Items {
		transporterPerson := commonModels.TransporterContactPersonDto{}

		err = dynamodbattribute.UnmarshalMap(val, &transporterPerson)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		transporterPersons[i] = transporterPerson
	}
	return transporterPersons, nil
}

func (repo *TransporterPersistance) GetTransporter(request commonModels.GetTransporterRequestDto) (commonModels.TransporterDto, *commonModels.ErrorDetail) {

	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").Equal(expression.Value(common.GetTransporterSortKey(request.TransporterId))),
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
		TableName:                 aws.String(repo.transporterTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get transporter call failed: %s", err.Error()))
	}

	transporter := commonModels.TransporterDto{}
	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &transporter)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
	}
	return transporter, nil
}

func buildFilterExpression(filterData commonModels.TransporterListRequest, projection *expression.ProjectionBuilder) (*expression.Expression, *commonModels.ErrorDetail) {

	filter := expression.Name("branchId").Equal(expression.Value(filterData.BranchId)).And(expression.Name("sortKey").BeginsWith(common.TransporterSortKey))

	if len(filterData.Alias) > 0 {
		filter = filter.And(expression.Name("alias").Contains(filterData.Alias).Or(expression.Name("alias").Equal(expression.Value(filterData.Alias))).Or(expression.Name("alias").BeginsWith(filterData.Alias)))
	}

	if len(filterData.CompanyName) > 0 {
		filter = filter.And(expression.Name("companyName").Contains(filterData.CompanyName).Or(expression.Name("companyName").Equal(expression.Value(filterData.CompanyName))).Or(expression.Name("companyName").BeginsWith(filterData.CompanyName)))
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

func (repo *TransporterPersistance) GetTransporterTotalByFilter(filterData commonModels.TransporterListRequest) (int64, *commonModels.ErrorDetail) {
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
		TableName:                 aws.String(repo.transporterTableName),
		ExclusiveStartKey:         filterData.LastEvalutionKey,
		Limit:                     aws.Int64(filterData.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter transporter call failed: %s", err.Error()))
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))
	for len(result.Items) > 0 && result.LastEvaluatedKey != nil {
		result, err = repo.db.Scan(&dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 aws.String(repo.transporterTableName),
			ExclusiveStartKey:         filterData.LastEvalutionKey,
			Limit:                     aws.Int64(filterData.PageSize),
		})
		if err != nil {
			common.WriteLog(1, fmt.Sprintf("filter transporter call failed: %s", err.Error()))
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}
	return count, nil
}

func (repo *TransporterPersistance) scanTransporter(expr *expression.Expression, filterdata *commonModels.TransporterListRequest) (*dynamodb.ScanOutput, *commonModels.ErrorDetail) {
	result, err := repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.transporterTableName),
		ExclusiveStartKey:         filterdata.LastEvalutionKey,
		Limit:                     aws.Int64(filterdata.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter transporter call failed: %s", err.Error()))
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: err.Error(),
		}
	}
	return result, nil
}

func parseTransporterScanResult(reultItem []map[string]*dynamodb.AttributeValue) []commonModels.TransporterDto {
	transporters := make([]commonModels.TransporterDto, 0)
	if len(reultItem) > 0 {
		for _, val := range reultItem {
			transporter := commonModels.TransporterDto{}

			err := dynamodbattribute.UnmarshalMap(val, &transporter)

			if err != nil {
				log.Fatalf("Got error unmarshalling: %s", err)
			}
			transporters = append(transporters, transporter)
		}
	}
	return transporters
}

func (repo *TransporterPersistance) GetTransporterByFilter(filterData commonModels.TransporterListRequest) ([]commonModels.TransporterDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {

	expr, errorDetails := buildFilterExpression(filterData, nil)
	if errorDetails != nil {
		return nil, nil, errorDetails
	}

	result, errorDetails := repo.scanTransporter(expr, &filterData)

	if errorDetails != nil {
		return nil, nil, errorDetails
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	transporters := parseTransporterScanResult(result.Items)

	for len(transporters) < int(filterData.PageSize) && filterData.LastEvalutionKey != nil {

		result, errorDetails = repo.scanTransporter(expr, &filterData)
		if errorDetails != nil {
			return nil, nil, errorDetails
		}

		transporterTemp := parseTransporterScanResult(result.Items)
		if len(transporterTemp) > 0 {
			transporters = append(transporters, transporterTemp...)
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
	}

	return transporters, filterData.LastEvalutionKey, nil
}

func (repo *TransporterPersistance) UpsertTransporter(transporter commonModels.TransporterDto, isNew bool) (*commonModels.TransporterDto, *commonModels.ErrorDetail) {
	existigClients, _, errorDetails := repo.GetTransporterByFilter(commonModels.TransporterListRequest{
		TransporterFilterDto: commonModels.TransporterFilterDto{
			BranchId:    transporter.BranchId,
			CompanyName: transporter.CompanyName,
			Alias:       transporter.Alias,
			Email:       transporter.ContactInfo.Email,
		},
		PageSize: 10,
	})

	if errorDetails != nil {
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Error in validating exiting transporter, error: %s", errorDetails.Error()),
		}
	}
	if len(existigClients) > 0 {
		flag := true
		if !isNew {
			flag = false
			for _, val := range existigClients {
				if val.TransporterId != transporter.TransporterId {
					flag = true
				}
			}
		}
		if flag {
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorInsert,
				ErrorMessage: "similar transporter already exists",
			}
		}
	}
	av, err := dynamodbattribute.MarshalMap(transporter)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling transporter details item, transporter name - %s, transporter id - %s, err: %s", transporter.CompanyName, transporter.TransporterId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.transporterTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating transporter %s, transporter id %s, error message; %s", transporter.CompanyName, transporter.TransporterId, err.Error()),
		}
	}
	return &transporter, nil
}

func (repo *TransporterPersistance) UpsertTransporterContact(transporterContact commonModels.TransporterContactPersonDto) (*commonModels.TransporterContactPersonDto, *commonModels.ErrorDetail) {
	fmt.Println("adding contact - ", transporterContact)

	av, err := dynamodbattribute.MarshalMap(transporterContact)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new transporter Contact detailes clinet id: %s, err: %s", transporterContact.TransporterId, err.Error()),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.transporterTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding transporter contact (%s) for transporter id %s, error message; %s", transporterContact.FirstName, transporterContact.TransporterId, err.Error()),
		}
	}
	return &transporterContact, nil
}

func (repo *TransporterPersistance) DeleteTransporterContact(branchId, transporterId, contactId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.transporterTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetTransporterContactSortKey(transporterId, contactId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())

		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting transporter contact id (%s) for transporter id %s, error message; %s", contactId, transporterId, err.Error()),
		}
	}
	return nil
}

func (repo *TransporterPersistance) DeleteTransporter(branchId, transporterId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.transporterTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetTransporterSortKey(transporterId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting transporter id %s, error message; %s", transporterId, err.Error()),
		}
	}
	return nil
}

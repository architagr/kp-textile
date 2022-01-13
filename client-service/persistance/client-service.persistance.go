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
	db               *dynamodb.DynamoDB
	clientTableName  string
	contactTableName string
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
			db:               dynamodb.New(dynamoDbSession),
			clientTableName:  common.EnvValues.ClientTableName,
			contactTableName: common.EnvValues.ContactTableName,
		}
	}

	return clientServicePersistanceObj, nil
}

func (repo *ClientServicePersistance) GetClientByFilter(filterData commonModels.ClientFilterDto) ([]commonModels.ClientDto, *commonModels.ErrorDetail) {

	filter := expression.Name("branchId").Contains(filterData.BranchId)

	if len(filterData.CompanyName) > 0 {
		filter = filter.And(expression.Name("compnayName").Contains(filterData.CompanyName))
	}

	if len(filterData.Alias) > 0 {
		filter = filter.And(expression.Name("alias").Contains(filterData.Alias))
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

	expr, err := expression.NewBuilder().WithFilter(filter).Build()

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("Got error building expression: %s", err.Error()))
	}
	result, err := repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.clientTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter client call failed: %s", err.Error()))
	}
	clients := make([]commonModels.ClientDto, len(result.Items))
	for i, val := range result.Items {
		client := commonModels.ClientDto{}

		err = dynamodbattribute.UnmarshalMap(val, &client)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		clients[i] = client
	}
	return clients, nil
}

func (repo *ClientServicePersistance) AddClient(client commonModels.ClientDto) (*commonModels.ClientDto, *commonModels.ErrorDetail) {
	av, err := dynamodbattribute.MarshalMap(client)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new Client detailes item: %s", err),
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
			ErrorMessage: fmt.Sprintf("Error in adding Client %s, error message; %s", client.CompanyName, err.Error()),
		}
	}
	return &client, nil
}

func (repo *ClientServicePersistance) AddClientContact(clientContact commonModels.ContactPersonDto) (*commonModels.ContactPersonDto, *commonModels.ErrorDetail) {
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

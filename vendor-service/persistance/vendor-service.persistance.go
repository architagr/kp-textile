package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"log"
	"vendor-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var vendorPersistanceObj *VendorPersistance

type VendorPersistance struct {
	db              *dynamodb.DynamoDB
	vendorTableName string
}

func InitVendorPersistance() (*VendorPersistance, *commonModels.ErrorDetail) {
	if vendorPersistanceObj == nil {
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

		vendorPersistanceObj = &VendorPersistance{
			db:              dynamodb.New(dynamoDbSession),
			vendorTableName: common.EnvValues.VendorTableName,
		}
	}

	return vendorPersistanceObj, nil
}

func (repo *VendorPersistance) GetPersonByVendorId(request commonModels.GetVendorRequestDto) ([]commonModels.VendorContactPersonDto, *commonModels.ErrorDetail) {
	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").BeginsWith(common.GetVendorContactSortKey(request.VendorId, "")),
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
		TableName:                 aws.String(repo.vendorTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get person for vendor %s call failed: %s", request.VendorId, err.Error()))
	}

	vendorPersons := make([]commonModels.VendorContactPersonDto, len(result.Items))

	for i, val := range result.Items {
		vendorPerson := commonModels.VendorContactPersonDto{}

		err = dynamodbattribute.UnmarshalMap(val, &vendorPerson)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		vendorPersons[i] = vendorPerson
	}
	return vendorPersons, nil
}

func (repo *VendorPersistance) GetVendor(request commonModels.GetVendorRequestDto) (commonModels.VendorDto, *commonModels.ErrorDetail) {

	keyCondition := expression.KeyAnd(
		expression.Key("branchId").Equal(expression.Value(request.BranchId)),
		expression.Key("sortKey").Equal(expression.Value(common.GetVendorSortKey(request.VendorId))),
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
		TableName:                 aws.String(repo.vendorTableName),
	})

	if err != nil {
		common.WriteLog(1, fmt.Sprintf("get vendor call failed: %s", err.Error()))
	}

	vendor := commonModels.VendorDto{}
	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &vendor)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
	}
	return vendor, nil
}

func buildFilterExpression(filterData commonModels.VendorListRequest, projection *expression.ProjectionBuilder) (*expression.Expression, *commonModels.ErrorDetail) {

	filter := expression.Name("branchId").Equal(expression.Value(filterData.BranchId)).And(expression.Name("sortKey").BeginsWith(common.VendorSortKey))

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

func (repo *VendorPersistance) GetVendorTotalByFilter(filterData commonModels.VendorListRequest) (int64, *commonModels.ErrorDetail) {
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
		TableName:                 aws.String(repo.vendorTableName),
		ExclusiveStartKey:         filterData.LastEvalutionKey,
		Limit:                     aws.Int64(filterData.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter vendor call failed: %s", err.Error()))
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	count = count + int64(len(result.Items))
	for len(result.Items) > 0 && result.LastEvaluatedKey != nil {
		result, err = repo.db.Scan(&dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 aws.String(repo.vendorTableName),
			ExclusiveStartKey:         filterData.LastEvalutionKey,
			Limit:                     aws.Int64(filterData.PageSize),
		})
		if err != nil {
			common.WriteLog(1, fmt.Sprintf("filter vendor call failed: %s", err.Error()))
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
		count = count + int64(len(result.Items))
	}
	return count, nil
}

func (repo *VendorPersistance) scanVendor(expr *expression.Expression, filterdata *commonModels.VendorListRequest) (*dynamodb.ScanOutput, *commonModels.ErrorDetail) {
	result, err := repo.db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.vendorTableName),
		ExclusiveStartKey:         filterdata.LastEvalutionKey,
		Limit:                     aws.Int64(filterdata.PageSize),
	})
	if err != nil {
		common.WriteLog(1, fmt.Sprintf("filter vendor call failed: %s", err.Error()))
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: err.Error(),
		}
	}
	return result, nil
}

func parseVendorScanResult(reultItem []map[string]*dynamodb.AttributeValue) []commonModels.VendorDto {
	vendors := make([]commonModels.VendorDto, 0)
	if len(reultItem) > 0 {
		for _, val := range reultItem {
			vendor := commonModels.VendorDto{}

			err := dynamodbattribute.UnmarshalMap(val, &vendor)

			if err != nil {
				log.Fatalf("Got error unmarshalling: %s", err)
			}
			vendors = append(vendors, vendor)
		}
	}
	return vendors
}

func (repo *VendorPersistance) GetVendorByFilter(filterData commonModels.VendorListRequest) ([]commonModels.VendorDto, map[string]*dynamodb.AttributeValue, *commonModels.ErrorDetail) {

	expr, errorDetails := buildFilterExpression(filterData, nil)
	if errorDetails != nil {
		return nil, nil, errorDetails
	}

	result, errorDetails := repo.scanVendor(expr, &filterData)

	if errorDetails != nil {
		return nil, nil, errorDetails
	}
	filterData.LastEvalutionKey = result.LastEvaluatedKey
	vendors := parseVendorScanResult(result.Items)

	for len(vendors) < int(filterData.PageSize) && filterData.LastEvalutionKey != nil {

		result, errorDetails = repo.scanVendor(expr, &filterData)
		if errorDetails != nil {
			return nil, nil, errorDetails
		}

		vendorsTemp := parseVendorScanResult(result.Items)
		if len(vendorsTemp) > 0 {
			vendors = append(vendors, vendorsTemp...)
		}
		filterData.LastEvalutionKey = result.LastEvaluatedKey
	}

	return vendors, filterData.LastEvalutionKey, nil
}

func (repo *VendorPersistance) UpsertVendor(vendor commonModels.VendorDto, isNew bool) (*commonModels.VendorDto, *commonModels.ErrorDetail) {
	existigClients, _, errorDetails := repo.GetVendorByFilter(commonModels.VendorListRequest{
		VendorFilterDto: commonModels.VendorFilterDto{
			BranchId:    vendor.BranchId,
			CompanyName: vendor.CompanyName,
			Alias:       vendor.Alias,
			Email:       vendor.ContactInfo.Email,
		},
		PageSize: 10,
	})

	if errorDetails != nil {
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Error in validating exiting vendor, error: %s", errorDetails.Error()),
		}
	}
	if len(existigClients) > 0 {
		flag := true
		if !isNew {
			flag = false
			for _, val := range existigClients {
				if val.VendorId != vendor.VendorId {
					flag = true
				}
			}
		}
		if flag {
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorInsert,
				ErrorMessage: "similar vendor already exists",
			}
		}
	}
	av, err := dynamodbattribute.MarshalMap(vendor)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling vendor details item, vendor name - %s, vendor id - %s, err: %s", vendor.CompanyName, vendor.VendorId, err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.vendorTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding/updating vendor %s, vendor id %s, error message; %s", vendor.CompanyName, vendor.VendorId, err.Error()),
		}
	}
	return &vendor, nil
}

func (repo *VendorPersistance) UpsertVendorContact(vendorContact commonModels.VendorContactPersonDto) (*commonModels.VendorContactPersonDto, *commonModels.ErrorDetail) {
	fmt.Println("adding contact - ", vendorContact)

	av, err := dynamodbattribute.MarshalMap(vendorContact)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new vendor Contact detailes clinet id: %s, err: %s", vendorContact.VendorId, err.Error()),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.vendorTableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding vendor contact (%s) for vendor id %s, error message; %s", vendorContact.FirstName, vendorContact.VendorId, err.Error()),
		}
	}
	return &vendorContact, nil
}

func (repo *VendorPersistance) DeleteVendorContact(branchId, vendorId, contactId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.vendorTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetVendorContactSortKey(vendorId, contactId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())

		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting vendor contact id (%s) for vendor id %s, error message; %s", contactId, vendorId, err.Error()),
		}
	}
	return nil
}

func (repo *VendorPersistance) DeleteVendor(branchId, vendorId string) *commonModels.ErrorDetail {
	_, err := repo.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &repo.vendorTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"branchId": {
				S: aws.String(branchId),
			},
			"sortKey": {
				S: aws.String(common.GetVendorSortKey(vendorId)),
			},
		},
	})
	if err != nil {
		common.WriteLog(1, err.Error())
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorDelete,
			ErrorMessage: fmt.Sprintf("Error in deleting vendor id %s, error message; %s", vendorId, err.Error()),
		}
	}
	return nil
}

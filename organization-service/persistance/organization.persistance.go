package persistance

import (
	commonModels "commonpkg/models"
	"encoding/json"
	"fmt"
	"organization-service/common"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/iris-contrib/go.uuid"
)

var organizationpersistanceObj *OrganizationPersistance

type getItemsResponseHeaders struct {
	ContentType []string `json:"Content-Type"`
}

type getItemsResponse struct {
	StatusCode int                     `json:"statusCode"`
	Headers    getItemsResponseHeaders `json:"multiValueHeaders"`
	Body       string                  `json:"body"`
}

type OrganizationPersistance struct {
	db        *dynamodb.DynamoDB
	lambda    *lambda.Lambda
	tableName string
}

func InitOrganizationPersistance() (*OrganizationPersistance, *commonModels.ErrorDetail) {
	if organizationpersistanceObj == nil {
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

		organizationpersistanceObj = &OrganizationPersistance{
			db:        dynamodb.New(dynamoDbSession),
			lambda:    lambda.New(dynamoDbSession),
			tableName: common.EnvValues.GodownTableName,
		}
	}

	return organizationpersistanceObj, nil
}

func (repo *OrganizationPersistance) GetAll() ([]commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	invocationType := "RequestResponse"
	functionName := "arn:aws:lambda:ap-south-1:675174225340:function:hsn-code-int-lambda-fn"
	var request = events.APIGatewayProxyRequest{
		Path:       "/",
		HTTPMethod: "get",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling MyGetItemsFunction request")
	}

	result, err := repo.lambda.Invoke(&lambda.InvokeInput{FunctionName: &functionName, InvocationType: &invocationType, Payload: payload})

	if err != nil {
		fmt.Println("Error calling MyGetItemsFunction")
	}
	var resp getItemsResponse
	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling MyGetItemsFunction response", err)
	}
	var body commonModels.HnsCodeListResponse

	err = json.Unmarshal([]byte(resp.Body), &body)
	if err != nil {
		fmt.Println("Error unmarshalling body response", err)
	}

	if body.StatusCode != 200 {
		fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(body.StatusCode))
	}

	// Print out items
	if body.Total > 0 {
		fmt.Printf("data received %+v\n", body.Data)
	} else {
		fmt.Println("There were no items")
	}

	return body.Data, nil
}

func (repo *OrganizationPersistance) Get(id string) (*commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	var hnsCode *commonModels.HnsCodeDto
	result, err := repo.db.Query(&dynamodb.QueryInput{
		TableName: &repo.tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
		},
	})

	if err != nil {
		common.WriteLog(1, err.Error())
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: err.Error(),
		}
	}
	if result.Items == nil || len(result.Items) == 0 {
		message := fmt.Sprintf("Could not find HSN codes for id %s", id)
		common.WriteLog(3, message)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}
	if len(result.Items) > 0 {

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &hnsCode)
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: "hsn codes not is correct format",
			}
		}
	}
	return hnsCode, nil
}

func (repo *OrganizationPersistance) Add(code string) (*commonModels.HnsCodeDto, *commonModels.ErrorDetail) {
	id, _ := uuid.NewV1()
	newHnsCode := commonModels.HnsCodeDto{Id: id.String(), HnsCode: code}

	av, err := dynamodbattribute.MarshalMap(newHnsCode)
	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("Got error marshalling new HnsCode item: %s", err),
		}
	}
	_, err = repo.db.PutItem(&dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})

	if err != nil {
		common.WriteLog(1, err.Error())

		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorInsert,
			ErrorMessage: fmt.Sprintf("Error in adding HSN code %s, error message; %s", code, err.Error()),
		}
	}

	return &newHnsCode, nil
}

func (repo *OrganizationPersistance) AddMultiple(codes []string) ([]commonModels.HnsCodeDto, []commonModels.ErrorDetail) {
	var newHnsCodes []commonModels.HnsCodeDto
	var errors = make([]commonModels.ErrorDetail, 0)
	for _, val := range codes {
		newHnsCode, err := repo.Add(val)
		if err != nil {
			common.WriteLog(1, err.Error())
			errors = append(errors, *err)
		} else {
			newHnsCodes = append(newHnsCodes, *newHnsCode)
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	return newHnsCodes, nil
}

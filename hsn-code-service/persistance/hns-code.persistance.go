package persistance

import (
	"fmt"
	"hsn-code-service/model"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var hnsCodes = []model.HnsCodeDto{
	{
		Id:      "1",
		HnsCode: "1-123",
	},
	{
		Id:      "2",
		HnsCode: "2-123",
	},
}
var hnsCodepersistanceObj *HnsCodePersistance

type HnsCodePersistance struct {
	db *dynamodb.DynamoDB
}

func InitHnsCodePersistance() *HnsCodePersistance {
	// TODO: get connection to DB
	if hnsCodepersistanceObj == nil {
		dynamoDbSession := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		hnsCodepersistanceObj = &HnsCodePersistance{
			db: dynamodb.New(dynamoDbSession),
		}
	}

	return hnsCodepersistanceObj
}

func (repo *HnsCodePersistance) GetAll() []model.HnsCodeDto {

	result, err := repo.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String("hsn-code"),
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Items == nil {
		msg := "Could not find hsn codes"
		log.Fatalf(msg)
	}
	items := make([]model.HnsCodeDto, 0)

	for _, val := range result.Items {
		item := model.HnsCodeDto{}
		err = dynamodbattribute.UnmarshalMap(val, &item)
		if err != nil {
			log.Fatal("hsn codes not is correct format")
		}
		items = append(items, item)
	}
	return items
}

func (repo *HnsCodePersistance) Get(id string) model.HnsCodeDto {
	var hnsCode *model.HnsCodeDto
	result, err := repo.db.Query(&dynamodb.QueryInput{
		TableName: aws.String("hsn-code"),
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
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Items == nil {
		msg := fmt.Sprintf("Could not find hsn codes for id %s", id)
		log.Fatalf(msg)
	}
	if len(result.Items) > 0 {

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &hnsCode)
		if err != nil {
			log.Fatal("hsn codes not is correct format")
		}
	} else {
		log.Fatal("hsn codes not found")
	}
	return *hnsCode
}

func (repo *HnsCodePersistance) Add(code string) model.HnsCodeDto {
	length := len(hnsCodes)

	id := hnsCodes[length-1].Id + "1"

	newHnsCode := model.HnsCodeDto{Id: id, HnsCode: code}
	hnsCodes = append(hnsCodes, newHnsCode)
	return newHnsCode
}

func (repo *HnsCodePersistance) AddMultiple(codes []string) []model.HnsCodeDto {
	var newHnsCodes []model.HnsCodeDto
	for _, val := range codes {
		newHnsCode := repo.Add(val)

		newHnsCodes = append(newHnsCodes, newHnsCode)
	}

	// TODO: return error of the codes which have not been added
	return newHnsCodes
}

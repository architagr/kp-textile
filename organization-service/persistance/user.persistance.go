package persistance

import (
	commonModels "commonpkg/models"
	"fmt"
	"organization-service/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var userPersistanceObj *UserPersistance

type UserPersistance struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func InitUserPersistance() (*UserPersistance, *commonModels.ErrorDetail) {
	if userPersistanceObj == nil {
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

		userPersistanceObj = &UserPersistance{
			db:        dynamodb.New(dynamoDbSession),
			tableName: common.EnvValues.UserTableName,
		}
	}

	return userPersistanceObj, nil
}

func (repo *UserPersistance) GetUserByUserName(userName string) (*commonModels.UserDto, *commonModels.ErrorDetail) {
	var user *commonModels.UserDto
	result, err := repo.db.Query(&dynamodb.QueryInput{
		TableName: &repo.tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"userName": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userName),
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
		message := fmt.Sprintf("Could not find user for username: %s", userName)
		common.WriteLog(3, message)
		return nil, &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorNoDataFound,
			ErrorMessage: message,
		}
	}
	if len(result.Items) > 0 {

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
		if err != nil {
			common.WriteLog(1, err.Error())
			return nil, &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorNoDataFound,
				ErrorMessage: "username not is correct format",
			}
		}
	}
	return user, nil
}

package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	// db error codes start from 10--
	ErrorDbConnection = "KP1001"
	ErrorNoDataFound  = "KP1002"
	ErrorInsert       = "KP1003"
	ErrorUpdate       = "KP1004"
	ErrorDelete       = "KP1005"

	// request param validation error codes start from 11--
	ErrorNoRequestData              = "KP1101"
	ErrorMissingRequiredRequestData = "KP1102"
	ErrorNoRequestParam             = "KP1103"
	ErrorInvalidRequestParam        = "KP1104"

	// general error codes start from 00--
	ErrorServer = "KP0001"
)

type ErrorDetail struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (err *ErrorDetail) Error() string {
	return fmt.Sprintf("ErrorCode: %s : ErrorMessage: %s", err.ErrorCode, err.ErrorMessage)
}

type CommonResponse struct {
	StatusCode   int           `json:"statusCode,omitempty"`
	ErrorMessage string        `json:"errorMessage,omitempty"`
	Errors       []ErrorDetail `json:"errors,omitempty"`
}
type CommonListResponse struct {
	CommonResponse
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty"`
	Total            int64                               `json:"total,omitempty"`
}

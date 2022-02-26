package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type TransporterDto struct {
	TransporterId         string            `json:"transporterId,omitempty" uri:"transporterId"`
	SortKey               string            `json:"sortKey,omitempty"`
	CompanyName           string            `json:"companyName,omitempty"`
	Alias                 string            `json:"alias,omitempty"`
	ContactInfo           ContactDetailsDto `json:"contactInfo,omitempty"`
	Website               string            `json:"website,omitempty"`
	PaymentTerm           string            `json:"paymentTerm,omitempty"`
	Addresses             []AddressDto      `json:"addresses,omitempty"`
	Remark                string            `json:"remark,omitempty"`
	Status                string            `json:"status,omitempty"`
	Gstn                  string            `json:"gstn,omitempty"`
	ParentTransporterId   string            `json:"parentTransporterId,omitempty"`
	ParentTransporterPath string            `json:"parentTransporterPath,omitempty"`
}

type TransporterContactPersonDto struct {
	TransporterId string            `json:"transporterId,omitempty"`
	SortKey       string            `json:"sortKey,omitempty"`
	ContactId     string            `json:"contactId,omitempty"`
	Salutation    string            `json:"salutation,omitempty"`
	FirstName     string            `json:"firstName,omitempty"`
	LastName      string            `json:"lastName,omitempty"`
	ContactInfo   ContactDetailsDto `json:"contactInfo,omitempty"`
	Department    string            `json:"department,omitempty"`
	PersonType    string            `json:"personType,omitempty"`
	Remark        string            `json:"remark,omitempty"`
	Address       AddressDto        `json:"address,omitempty"`
}

type AddTransporterRequest struct {
	TransporterDto
	ContactPersons []TransporterContactPersonDto `json:"contactPersons,omitempty"`
}

type AddTransporterResponse struct {
	CommonResponse
	Data AddTransporterRequest `json:"data,omitempty"`
}

type TransporterFilterDto struct {
	CompanyName            string `json:"companyName,omitempty" form:"companyName"`
	Alias                  string `json:"alias,omitempty" form:"alias"`
	Email                  string `json:"email,omitempty" form:"email"`
	ContactPersonFirstName string `json:"contactPersonFirstName,omitempty" form:"contactPersonFirstName"`
	ContactPersonLastName  string `json:"contactPersonLastName,omitempty" form:"contactPersonLastName"`
	PaymentTerm            string `json:"paymentTerm,omitempty" form:"paymentTerm"`
}

type TransporterListRequest struct {
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty" form:"pageSize"`
	TransporterFilterDto
}

type TransporterListResponse struct {
	CommonListResponse
	Data []TransporterDto `json:"data,omitempty"`
}

type TransporterRequest struct {
	Id string `uri:"id"`
}

type TransporterResponse struct {
	CommonResponse
	Data TransporterDto `json:"data,omitempty"`
}

type GetTransporterRequestDto struct {
	TransporterId string `uri:"transporterId"`
}

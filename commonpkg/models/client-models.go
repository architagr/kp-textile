package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	Status_Active   = "Active"
	Status_Blocked  = "Blocked"
	Status_Inactive = "Inactive"
)

const (
	PaymentTerm_15Days  = "15Days"
	PaymentTerm_Monthly = "Monthly"
	PaymentTerm_Bill    = "Bill"
)

const (
	AddressType_HeadOffice     = "HeadOffice"
	AddressType_SalesOffice    = "SalesOffice"
	AddressType_RegionalOffice = "RegionalOffice"
	AddressType_ShippingOffice = "ShippingOffice"
	AddressType_Residence      = "Residence"
	AddressType_Warehouse      = "Warehouse"
)

const (
	PersonType_Sales     = "Sales"
	PersonType_Account   = "Account"
	PersonType_Store     = "Store"
	PersonType_Marketing = "Marketing"
	PersonType_Director  = "Director"
	PersonType_Chairman  = "Chairman"
)

type ClientDto struct {
	BranchId         string            `json:"branchId,omitempty"`
	SortKey          string            `json:"sortKey,omitempty"`
	ClientId         string            `json:"clientId,omitempty"`
	CompanyName      string            `json:"companyName,omitempty"`
	Alias            string            `json:"alias,omitempty"`
	ContactInfo      ContactDetailsDto `json:"contactInfo,omitempty"`
	Website          string            `json:"website,omitempty"`
	PaymentTerm      string            `json:"paymentTerm,omitempty"`
	Addresses        []AddressDto      `json:"addresses,omitempty"`
	Remark           string            `json:"remark,omitempty"`
	Status           string            `json:"status,omitempty"`
	Gstn             string            `json:"gstn,omitempty"`
	ParentClientId   string            `json:"parentClientId,omitempty"`
	ParentClientPath string            `json:"parentClientPath,omitempty"`
}
type ContactPersonDto struct {
	BranchId    string            `json:"branchId,omitempty"`
	SortKey     string            `json:"sortKey,omitempty"`
	ClientId    string            `json:"clientId,omitempty"`
	ContactId   string            `json:"contactId,omitempty"`
	Salutation  string            `json:"salutation,omitempty"`
	FirstName   string            `json:"firstName,omitempty"`
	LastName    string            `json:"lastName,omitempty"`
	ContactInfo ContactDetailsDto `json:"contactInfo,omitempty"`
	Department  string            `json:"department,omitempty"`
	PersonType  string            `json:"personType,omitempty"`
	Remark      string            `json:"remark,omitempty"`
	Address     AddressDto        `json:"address,omitempty"`
}

type AddClientRequest struct {
	ClientDto
	ContactPersons []ContactPersonDto `json:"contactPersons,omitempty"`
}

type AddClientResponse struct {
	CommonResponse
	Data AddClientRequest `json:"data,omitempty"`
}

type ContactDetailsDto struct {
	Email    string `json:"email,omitempty"`
	Landline string `json:"landline,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Whatsapp string `json:"whatsapp,omitempty"`
}

type AddressDto struct {
	Label        string `json:"label,omitempty"`
	AddressType  string `json:"addressType,omitempty"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	Pincode      string `json:"pincode,omitempty"`
	Country      string `json:"country,omitempty"`
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	Zipcode      string `json:"zipcode,omitempty"`
	Landline     string `json:"landline,omitempty"`
	Mobile       string `json:"mobile,omitempty"`
}

type ClientFilterDto struct {
	BranchId               string
	CompanyName            string `json:"companyName,omitempty" form:"companyName"`
	Alias                  string `json:"alias,omitempty" form:"alias"`
	Email                  string `json:"email,omitempty" form:"email"`
	ContactPersonFirstName string `json:"contactPersonFirstName,omitempty" form:"contactPersonFirstName"`
	ContactPersonLastName  string `json:"contactPersonLastName,omitempty" form:"contactPersonLastName"`
	PaymentTerm            string `json:"paymentTerm,omitempty" form:"paymentTerm"`
}

type ClientListRequest struct {
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty" form:"pageSize"`
	ClientFilterDto
}

type ClientListResponse struct {
	CommonListResponse
	Data []ClientDto `json:"data,omitempty"`
}

type ClientRequest struct {
	Id string `uri:"id"`
}

type ClientResponse struct {
	CommonResponse
	Data ClientDto `json:"data,omitempty"`
}

type GetClientRequestDto struct {
	BranchId string
	ClientId string `uri:"clientId"`
}

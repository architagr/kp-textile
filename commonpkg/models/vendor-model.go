package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type VendorDto struct {
	BranchId         string            `json:"branchId,omitempty"`
	SortKey          string            `json:"sortKey,omitempty"`
	VendorId         string            `json:"vendorId,omitempty" uri:"vendorId"`
	CompanyName      string            `json:"companyName,omitempty"`
	Alias            string            `json:"alias,omitempty"`
	ContactInfo      ContactDetailsDto `json:"contactInfo,omitempty"`
	Website          string            `json:"website,omitempty"`
	PaymentTerm      string            `json:"paymentTerm,omitempty"`
	Addresses        []AddressDto      `json:"addresses,omitempty"`
	Remark           string            `json:"remark,omitempty"`
	Status           string            `json:"status,omitempty"`
	Gstn             string            `json:"gstn,omitempty"`
	ParentVendorId   string            `json:"parentVendorId,omitempty"`
	ParentVendorPath string            `json:"parentVendorPath,omitempty"`
}

type VendorContactPersonDto struct {
	BranchId    string            `json:"branchId,omitempty"`
	SortKey     string            `json:"sortKey,omitempty"`
	VendorId    string            `json:"vendorId,omitempty"`
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

type AddVendorRequest struct {
	VendorDto
	ContactPersons []VendorContactPersonDto `json:"contactPersons,omitempty"`
}

type AddVendorResponse struct {
	CommonResponse
	Data AddVendorRequest `json:"data,omitempty"`
}

type VendorFilterDto struct {
	BranchId               string
	CompanyName            string `json:"companyName,omitempty" form:"companyName"`
	Alias                  string `json:"alias,omitempty" form:"alias"`
	Email                  string `json:"email,omitempty" form:"email"`
	ContactPersonFirstName string `json:"contactPersonFirstName,omitempty" form:"contactPersonFirstName"`
	ContactPersonLastName  string `json:"contactPersonLastName,omitempty" form:"contactPersonLastName"`
	PaymentTerm            string `json:"paymentTerm,omitempty" form:"paymentTerm"`
}

type VendorListRequest struct {
	LastEvalutionKey map[string]*dynamodb.AttributeValue `json:"lastEvalutionKey,omitempty"`
	PageSize         int64                               `json:"pageSize,omitempty" form:"pageSize"`
	VendorFilterDto
}

type VendorListResponse struct {
	CommonListResponse
	Data []VendorDto `json:"data,omitempty"`
}

type VendorRequest struct {
	Id string `uri:"id"`
}

type VendorResponse struct {
	CommonResponse
	Data VendorDto `json:"data,omitempty"`
}

type GetVendorRequestDto struct {
	BranchId string
	VendorId string `uri:"vendorId"`
}

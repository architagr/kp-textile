package models

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
	CompanyName            string `json:"companyName,omitempty"`
	Alias                  string `json:"alias,omitempty"`
	Email                  string `json:"email,omitempty"`
	ContactPersonFirstName string `json:"contactPersonFirstName,omitempty"`
	ContactPersonLastName  string `json:"contactPersonLastName,omitempty"`
	PaymentTerm            string `json:"paymentTerm,omitempty"`
}

type ClientListRequest struct {
	Start    int             `json:"start,omitempty"`
	PageSize int             `json:"pageSize,omitempty"`
	Filter   ClientFilterDto `json:"filter,omitempty"`
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

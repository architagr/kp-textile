import { CommonListResponse, CommonResponse } from "./genric-model";

export interface ClientDto {
    branchId: string
    sortKey: string
    clientId: string
    companyName: string
    alias: string
    contactInfo: ContactDetailsDto
    website: string
    paymentTerm: string
    addresses: AddressDto[]
    remark: string
    status: string
    parentClientId: string
    parentClientPath: string;
    gstn: string;
}

export interface ContactDetailsDto {
    email: string
    landline: string
    mobile: string
    whatsapp: string
}


export interface AddressDto {
    label: string
    addressType: string
    addressLine1: string
    addressLine2: string
    pincode: string
    country: string
    state: string
    city: string
    landline: string
    mobile: string
}

export interface ContactPersonDto {
    branchId: string
    sortKey: string
    clientId: string
    contactId: string
    salutation: string
    firstName: string
    lastName: string
    contactInfo: ContactDetailsDto
    department: string
    personType: string
    remark: string
    address: AddressDto
}

export interface ClientListResponse extends CommonListResponse {

    data: ClientDto[]
}

export interface ClientResponse extends CommonResponse {
    data: ClientDto
}

export interface AddClientRequest extends ClientDto {

    contactPersons: ContactPersonDto[]
}

export interface AddClientResponse extends CommonResponse {
	data: AddClientRequest
}
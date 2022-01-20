import { AddressDto, ContactDetailsDto } from "./client-model";
import { CommonListResponse, CommonResponse } from "./genric-model"

export interface VendorDto {
    branchId: string;
    sortKey: string;
    vendorId: string;
    companyName: string;
    alias: string;
    contactInfo: ContactDetailsDto;
    website: string;
    paymentTerm: string;
    addresses: AddressDto[];
    remark: string;
    status: string;
    gstn: string;
    parentVendorId: string;
    parentVendorPath: string;
}

export interface VendorContactPersonDto {
    branchId: string;
    sortKey: string;
    vendorId: string;
    contactId: string;
    salutation: string;
    firstName: string;
    lastName: string;
    contactInfo: ContactDetailsDto;
    department: string;
    personType: string;
    remark: string;
    address: AddressDto;
}

export interface AddVendorRequest extends VendorDto {
    contactPersons: VendorContactPersonDto[];
}

export interface AddVendorResponse extends CommonResponse {
    data: AddVendorRequest
}

export interface VendorFilterDto {
    branchId: string;
    companyName: string;
    alias: string;
    email: string;
    contactPersonFirstName: string;
    contactPersonLastName: string;
    paymentTerm: string;
}

export interface VendorListRequest extends VendorFilterDto {
    lastEvalutionKey: any;
    pageSize: number;
}

export interface VendorListResponse extends CommonListResponse {
    data: VendorDto[];
}

export interface VendorRequest {
    id: string;
}

export interface VendorResponse extends CommonResponse {
    data: VendorDto
}

export interface GetVendorRequestDto {
    branchId: string;
    vendorId: string;
}

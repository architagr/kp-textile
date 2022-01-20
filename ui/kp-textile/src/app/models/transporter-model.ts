import { AddressDto, ContactDetailsDto } from "./client-model";
import { CommonListResponse, CommonResponse } from "./genric-model";

export interface TransporterDto {
    branchId: string;
    sortKey: string;
    transporterId: string;
    companyName: string;
    alias: string;
    contactInfo: ContactDetailsDto;
    website: string;
    paymentTerm: string;
    addresses: AddressDto[];
    remark: string;
    status: string;
    gstn: string;
    parentTransporterId: string;
    parentTransporterPath: string;
}

export interface TransporterContactPersonDto {
    branchId: string;
    sortKey: string;
    transporterId: string;
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

export interface AddTransporterRequest extends TransporterDto {
    contactPersons: TransporterContactPersonDto[]
}

export interface AddTransporterResponse extends CommonResponse {
    data: AddTransporterRequest
}

export interface TransporterFilterDto {
    branchId: string;
    companyName: string;
    alias: string;
    email: string;
    contactPersonFirstName: string;
    contactPersonLastName: string;
    paymentTerm: string;
}

export interface TransporterListRequest extends TransporterFilterDto {
    lastEvalutionKey: any
    pageSize: number

}

export interface TransporterListResponse extends CommonListResponse {
    data: TransporterDto[];
}

export interface TransporterRequest {
    id: string;
}

export interface TransporterResponse extends CommonResponse {
    data: TransporterDto;
}

export interface GetTransporterRequestDto {
    branchId: string;
    transporterId: string;
}

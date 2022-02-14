import { CommonListResponse, CommonResponse } from "./genric-model"

export interface GodownDto {
    id: string;
    name: string;
}

export interface GodownListResponse extends CommonListResponse {
    data: GodownDto[];
}

export interface GodownResponse extends CommonResponse {
    data: GodownDto;
}

export interface GodownAddRequest {
    name: string
}
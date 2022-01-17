import { CommonListResponse, CommonResponse } from "./genric-model"

export interface HnsCodeDto {
    id: string;
    hnsCode: string;
}


export interface  HnsCodeListResponse extends CommonListResponse {
	data: HnsCodeDto[];
}

export interface HnsCodeResponse extends CommonResponse {
	data: HnsCodeDto;
}


export interface AddRequest {
	code: string;
}
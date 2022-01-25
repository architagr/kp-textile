import { CommonListResponse, CommonResponse } from "./genric-model"

export interface QualityDto {
	id   :string
	name :string
}

export interface QualityListResponse extends CommonListResponse{
	data: QualityDto[]
}

export interface QualityResponse extends CommonResponse{
	data: QualityDto
}

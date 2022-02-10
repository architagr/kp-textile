import { CommonListResponse, CommonResponse } from "./genric-model"
import { BailDetailsDto } from "./item-model";

export interface QualityDto {
	id: string
	name: string
	hsnCode: string
	productId: string
	productName: string
}
export interface QualityListItemDto extends QualityDto {
	pendingQuantity: number
	bailDetails: BailDetailsDto[]
}
export interface QualityListResponse extends CommonListResponse {
	data: QualityDto[]
}

export interface QualityResponse extends CommonResponse {
	data: QualityDto
}

export interface ProductDto {
	id: string
	name: string
}

export interface ProductListResponse extends CommonListResponse {
	data: ProductDto[]
}

export interface ProductResponse extends CommonResponse {
	data: ProductDto
}
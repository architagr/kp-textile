package models

type QualityDto struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	HsnCode   string `json:"hsnCode"`
	ProductId string `json:"productId"`
}

type QualityListResponse struct {
	CommonListResponse
	Data []QualityDto `json:"data,omitempty"`
}

type QualityResponse struct {
	CommonResponse
	Data QualityDto `json:"data,omitempty"`
}

type ProductDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProductListResponse struct {
	CommonListResponse
	Data []ProductDto `json:"data,omitempty"`
}

type ProductResponse struct {
	CommonResponse
	Data ProductDto `json:"data,omitempty"`
}

package models

type HnsCodeDto struct {
	Id      string `json:"id"`
	HnsCode string `json:"hnsCode"`
}

type HnsCodeListResponse struct {
	CommonListResponse
	Data []HnsCodeDto `json:"data,omitempty"`
}

type HnsCodeResponse struct {
	CommonResponse
	Data HnsCodeDto `json:"data,omitempty"`
}

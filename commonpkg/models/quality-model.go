package models

type QualityDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type QualityListResponse struct {
	CommonListResponse
	Data []QualityDto `json:"data,omitempty"`
}

type QualityResponse struct {
	CommonResponse
	Data QualityDto `json:"data,omitempty"`
}

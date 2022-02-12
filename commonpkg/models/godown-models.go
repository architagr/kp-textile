package models

type GodownDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GodownListResponse struct {
	CommonListResponse
	Data []GodownDto `json:"data,omitempty"`
}

type GodownResponse struct {
	CommonResponse
	Data GodownDto `json:"data,omitempty"`
}

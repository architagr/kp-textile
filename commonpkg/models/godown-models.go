package models

type GodownDto struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GodownListResponse struct {
	CommonListResponse
	Data []GodownDto `json:"data,omitempty"`
}

type GodownResponse struct {
	CommonResponse
	Data GodownDto `json:"data,omitempty"`
}

type GodownAddRequest struct {
	Name string `json:"name"`
}

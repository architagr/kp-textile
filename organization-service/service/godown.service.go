package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"
	"organization-service/persistance"
)

var godownServiceObj *GodownService

type GodownService struct {
	repo *persistance.GodownPersistance
}

func InitGodownService() (*GodownService, *commonModels.ErrorDetail) {
	if godownServiceObj == nil {
		repo, err := persistance.InitGodownPersistance()
		if err != nil {
			return nil, err
		}
		godownServiceObj = &GodownService{
			repo: repo,
		}
	}
	return godownServiceObj, nil
}

func (svc *GodownService) GetAll() commonModels.GodownListResponse {
	data, err := svc.repo.GetAll()
	if err != nil {
		return commonModels.GodownListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting list of godowns",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	return commonModels.GodownListResponse{
		CommonListResponse: commonModels.CommonListResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			PageSize: int64(len(data)),
			Total:    int64(len(data)),
		},
		Data: data,
	}
}
func (svc *GodownService) Add(name string) commonModels.GodownResponse {
	response, err := svc.repo.Add(name)
	if err != nil {
		return commonModels.GodownResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Error in adding godowns with name %s", name),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	return commonModels.GodownResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: *response,
	}
}

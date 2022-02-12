package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"organization-service/persistance"
)

var organizationServiceObj *OrganizationService

type OrganizationService struct {
	organizationRepo *persistance.OrganizationPersistance
}

func InitHnsCodeService() (*OrganizationService, *commonModels.ErrorDetail) {
	if organizationServiceObj == nil {
		repo, err := persistance.InitOrganizationPersistance()
		if err != nil {
			return nil, err
		}
		organizationServiceObj = &OrganizationService{
			organizationRepo: repo,
		}
	}
	return organizationServiceObj, nil
}

func (service *OrganizationService) GetAll() commonModels.HnsCodeListResponse {
	allCodes, err := service.organizationRepo.GetAll()

	if err != nil {

		return commonModels.HnsCodeListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All HSN Codes",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.HnsCodeListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusOK,
				},
				LastEvalutionKey: nil,
				Total:            int64(len(allCodes)),
				PageSize:         int64(len(allCodes)),
			},
			Data: allCodes,
		}
	}
}

func (service *OrganizationService) Get(id string) commonModels.HnsCodeResponse {
	hsnCode, err := service.organizationRepo.Get(id)
	if err != nil {
		return commonModels.HnsCodeResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get HSN Code for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.HnsCodeResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *hsnCode,
		}
	}
}

func (service *OrganizationService) Add(code string) commonModels.HnsCodeResponse {

	hsnCode, err := service.organizationRepo.Add(code)

	if err != nil {
		return commonModels.HnsCodeResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add HSN Code - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.HnsCodeResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *hsnCode,
		}
	}
}

func (service *OrganizationService) AddMultiple(codes []string) commonModels.HnsCodeListResponse {
	allCodes, err := service.organizationRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.HnsCodeListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All HSN Codes",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.HnsCodeListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All HSN Codes",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.HnsCodeListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				LastEvalutionKey: nil,
				Total:            int64(len(allCodes)),
				PageSize:         int64(len(allCodes)),
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusCreated,
				},
			},
			Data: allCodes,
		}
	}
}

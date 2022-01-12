package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"transportor-service/persistance"
)

var TransportorServiceObj *TransportorServiceService

type TransportorServiceService struct {
	transportorServiceRepo *persistance.TransportorServicePersistance
}

func InitTransportorServiceService() (*TransportorServiceService, *commonModels.ErrorDetail) {
	if TransportorServiceObj == nil {
		repo, err := persistance.InitTransportorServicePersistance()
		if err != nil {
			return nil, err
		}
		TransportorServiceObj = &TransportorServiceService{
			transportorServiceRepo: repo,
		}
	}
	return TransportorServiceObj, nil
}

func (service *TransportorServiceService) GetAll() commonModels.TransportorServiceListResponse {
	allCodes, err := service.transportorServiceRepo.GetAll()

	if err != nil {

		return commonModels.TransportorServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All TransportorService",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.TransportorServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusOK,
				},
				Start:    0,
				Total:    len(allCodes),
				PageSize: len(allCodes),
			},
			Data: allCodes,
		}
	}
}

func (service *TransportorServiceService) Get(id string) commonModels.TransportorServiceResponse {
	transportorService, err := service.transportorServiceRepo.Get(id)
	if err != nil {
		return commonModels.TransportorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get HSN Code for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.TransportorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *transportorService,
		}
	}
}

func (service *TransportorServiceService) Add(code string) commonModels.TransportorServiceResponse {
	transportorService, err := service.transportorServiceRepo.Add(code)

	if err != nil {
		return commonModels.TransportorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add HSN Code - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.TransportorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *transportorService,
		}
	}
}

func (service *TransportorServiceService) AddMultiple(codes []string) commonModels.TransportorServiceListResponse {
	allCodes, err := service.transportorServiceRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.TransportorServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All TransportorService",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.TransportorServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All TransportorService",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.TransportorServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				Start:    0,
				Total:    len(allCodes),
				PageSize: len(allCodes),
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusCreated,
				},
			},
			Data: allCodes,
		}
	}
}

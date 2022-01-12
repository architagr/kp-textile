package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"item-service/persistance"
)

var ItemServiceObj *ItemServiceService

type ItemServiceService struct {
	itemServiceRepo *persistance.ItemServicePersistance
}

func InitItemServiceService() (*ItemServiceService, *commonModels.ErrorDetail) {
	if ItemServiceObj == nil {
		repo, err := persistance.InitItemServicePersistance()
		if err != nil {
			return nil, err
		}
		ItemServiceObj = &ItemServiceService{
			itemServiceRepo: repo,
		}
	}
	return ItemServiceObj, nil
}

func (service *ItemServiceService) GetAll() commonModels.ItemServiceListResponse {
	allCodes, err := service.itemServiceRepo.GetAll()

	if err != nil {

		return commonModels.ItemServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All ItemService",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.ItemServiceListResponse{
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

func (service *ItemServiceService) Get(id string) commonModels.ItemServiceResponse {
	itemService, err := service.itemServiceRepo.Get(id)
	if err != nil {
		return commonModels.ItemServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get HSN Code for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.ItemServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *itemService,
		}
	}
}

func (service *ItemServiceService) Add(code string) commonModels.ItemServiceResponse {
	itemService, err := service.itemServiceRepo.Add(code)

	if err != nil {
		return commonModels.ItemServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add HSN Code - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.ItemServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *itemService,
		}
	}
}

func (service *ItemServiceService) AddMultiple(codes []string) commonModels.ItemServiceListResponse {
	allCodes, err := service.itemServiceRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.ItemServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All ItemService",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.ItemServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All ItemService",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.ItemServiceListResponse{
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

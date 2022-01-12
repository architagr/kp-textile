package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"vendor-service/persistance"
)

var VendorServiceObj *VendorServiceService

type VendorServiceService struct {
	vendorServiceRepo *persistance.VendorServicePersistance
}

func InitVendorServiceService() (*VendorServiceService, *commonModels.ErrorDetail) {
	if VendorServiceObj == nil {
		repo, err := persistance.InitVendorServicePersistance()
		if err != nil {
			return nil, err
		}
		VendorServiceObj = &VendorServiceService{
			vendorServiceRepo: repo,
		}
	}
	return VendorServiceObj, nil
}

func (service *VendorServiceService) GetAll() commonModels.VendorServiceListResponse {
	allCodes, err := service.vendorServiceRepo.GetAll()

	if err != nil {

		return commonModels.VendorServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All VendorService",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.VendorServiceListResponse{
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

func (service *VendorServiceService) Get(id string) commonModels.VendorServiceResponse {
	vendorService, err := service.vendorServiceRepo.Get(id)
	if err != nil {
		return commonModels.VendorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get HSN Code for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.VendorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *vendorService,
		}
	}
}

func (service *VendorServiceService) Add(code string) commonModels.VendorServiceResponse {
	vendorService, err := service.vendorServiceRepo.Add(code)

	if err != nil {
		return commonModels.VendorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add HSN Code - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.VendorServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *vendorService,
		}
	}
}

func (service *VendorServiceService) AddMultiple(codes []string) commonModels.VendorServiceListResponse {
	allCodes, err := service.vendorServiceRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.VendorServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All VendorService",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.VendorServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All VendorService",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.VendorServiceListResponse{
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

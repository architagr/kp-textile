package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"hsn-code-service/persistance"
)

var hnsCodeServiceObj *HnsCodeService

type HnsCodeService struct {
	hnsCodeRepo *persistance.HnsCodePersistance
}

func InitHnsCodeService() (*HnsCodeService, *commonModels.ErrorDetail) {
	if hnsCodeServiceObj == nil {
		repo, err := persistance.InitHnsCodePersistance()
		if err != nil {
			return nil, err
		}
		hnsCodeServiceObj = &HnsCodeService{
			hnsCodeRepo: repo,
		}
	}
	return hnsCodeServiceObj, nil
}

func (service *HnsCodeService) GetAll() commonModels.HnsCodeListResponse {
	allCodes, err := service.hnsCodeRepo.GetAll()

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
				Total:    int64(len(allCodes)),
				PageSize: int64(len(allCodes)),
			},
			Data: allCodes,
		}
	}
}

func (service *HnsCodeService) Get(id string) commonModels.HnsCodeResponse {
	hsnCode, err := service.hnsCodeRepo.Get(id)
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

func (service *HnsCodeService) Add(code string) commonModels.HnsCodeResponse {
	fmt.Println("service", code)

	hsnCode, err := service.hnsCodeRepo.Add(code)

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

func (service *HnsCodeService) AddMultiple(codes []string) commonModels.HnsCodeListResponse {
	allCodes, err := service.hnsCodeRepo.AddMultiple(codes)

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
				Total:    int64(len(allCodes)),
				PageSize: int64(len(allCodes)),
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusCreated,
				},
			},
			Data: allCodes,
		}
	}
}

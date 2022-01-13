package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"quality-service/persistance"
)

var QualityServiceObj *QualityService

type QualityService struct {
	qualityRepo *persistance.QualityPersistance
}

func InitQualityService() (*QualityService, *commonModels.ErrorDetail) {
	if QualityServiceObj == nil {
		repo, err := persistance.InitQualityPersistance()
		if err != nil {
			return nil, err
		}
		QualityServiceObj = &QualityService{
			qualityRepo: repo,
		}
	}
	return QualityServiceObj, nil
}

func (service *QualityService) GetAll() commonModels.QualityListResponse {
	allCodes, err := service.qualityRepo.GetAll()

	if err != nil {

		return commonModels.QualityListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All QualityService",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.QualityListResponse{
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

func (service *QualityService) Get(id string) commonModels.QualityResponse {
	qualityService, err := service.qualityRepo.Get(id)
	if err != nil {
		return commonModels.QualityResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get Quality for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.QualityResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *qualityService,
		}
	}
}

func (service *QualityService) Add(code string) commonModels.QualityResponse {
	qualityService, err := service.qualityRepo.Add(code)

	if err != nil {
		return commonModels.QualityResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add Quality - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.QualityResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *qualityService,
		}
	}
}

func (service *QualityService) AddMultiple(codes []string) commonModels.QualityListResponse {
	allCodes, err := service.qualityRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.QualityListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All Qualities",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.QualityListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All Qualities",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.QualityListResponse{
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

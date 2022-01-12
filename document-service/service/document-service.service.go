package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"document-service/persistance"
)

var DocumentServiceObj *DocumentServiceService

type DocumentServiceService struct {
	documentServiceRepo *persistance.DocumentServicePersistance
}

func InitDocumentServiceService() (*DocumentServiceService, *commonModels.ErrorDetail) {
	if DocumentServiceObj == nil {
		repo, err := persistance.InitDocumentServicePersistance()
		if err != nil {
			return nil, err
		}
		DocumentServiceObj = &DocumentServiceService{
			documentServiceRepo: repo,
		}
	}
	return DocumentServiceObj, nil
}

func (service *DocumentServiceService) GetAll() commonModels.DocumentServiceListResponse {
	allCodes, err := service.documentServiceRepo.GetAll()

	if err != nil {

		return commonModels.DocumentServiceListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get All DocumentService",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	} else {
		return commonModels.DocumentServiceListResponse{
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

func (service *DocumentServiceService) Get(id string) commonModels.DocumentServiceResponse {
	documentService, err := service.documentServiceRepo.Get(id)
	if err != nil {
		return commonModels.DocumentServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get HSN Code for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.DocumentServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *documentService,
		}
	}
}

func (service *DocumentServiceService) Add(code string) commonModels.DocumentServiceResponse {
	documentService, err := service.documentServiceRepo.Add(code)

	if err != nil {
		return commonModels.DocumentServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add HSN Code - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.DocumentServiceResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *documentService,
		}
	}
}

func (service *DocumentServiceService) AddMultiple(codes []string) commonModels.DocumentServiceListResponse {
	allCodes, err := service.documentServiceRepo.AddMultiple(codes)

	if err != nil {

		if len(allCodes) > 0 && len(codes) > len(allCodes) {
			return commonModels.DocumentServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All DocumentService",
						Errors:       err,
					},
				},
				Data: allCodes,
			}
		} else {
			return commonModels.DocumentServiceListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusBadRequest,
						ErrorMessage: "could not add All DocumentService",
						Errors:       err,
					},
				},
			}
		}
	} else {
		return commonModels.DocumentServiceListResponse{
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

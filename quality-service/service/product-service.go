package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"quality-service/persistance"
)

var productServiceObj *ProductService

type ProductService struct {
	productRepo *persistance.ProductPersistance
}

func InitProductService() (*ProductService, *commonModels.ErrorDetail) {
	if productServiceObj == nil {
		repo, err := persistance.InitProductPersistance()
		if err != nil {
			return nil, err
		}
		productServiceObj = &ProductService{
			productRepo: repo,
		}
	}
	return productServiceObj, nil
}

func (service *ProductService) GetAll() commonModels.ProductListResponse {
	allProducts, err := service.productRepo.GetAll()
	if err != nil {

		return commonModels.ProductListResponse{
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
		return commonModels.ProductListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusOK,
				},
				Total:    int64(len(allProducts)),
				PageSize: int64(len(allProducts)),
			},
			Data: allProducts,
		}
	}
}

func (service *ProductService) Get(id string) commonModels.ProductResponse {
	product, err := service.productRepo.Get(id)
	if err != nil {
		return commonModels.ProductResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Could not get Quality for id: %s", id),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.ProductResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *product,
		}
	}
}

func (service *ProductService) Add(code string) commonModels.ProductResponse {
	product, err := service.productRepo.Add(code)

	if err != nil {
		return commonModels.ProductResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add Product - %s", code),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	} else {
		return commonModels.ProductResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusCreated,
			},
			Data: *product,
		}
	}
}

func (service *ProductService) AddMultiple(codes []string) commonModels.ProductListResponse {
	allProducts, err := service.productRepo.AddMultiple(codes)

	if err != nil {

		if len(allProducts) > 0 && len(codes) > len(allProducts) {
			return commonModels.ProductListResponse{
				CommonListResponse: commonModels.CommonListResponse{
					CommonResponse: commonModels.CommonResponse{
						StatusCode:   http.StatusPartialContent,
						ErrorMessage: "could not add All Qualities",
						Errors:       err,
					},
				},
				Data: allProducts,
			}
		} else {
			return commonModels.ProductListResponse{
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
		return commonModels.ProductListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				Total:    int64(len(allProducts)),
				PageSize: int64(len(allProducts)),
				CommonResponse: commonModels.CommonResponse{
					StatusCode: http.StatusCreated,
				},
			},
			Data: allProducts,
		}
	}
}

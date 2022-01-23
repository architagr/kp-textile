package service

import (
	commonModels "commonpkg/models"
	"item-service/persistance"
	"net/http"
)

var bailServiceObj *BailService

type BailService struct {
	bailRepo *persistance.BailPersistance
}

func InitBailService() (*BailService, *commonModels.ErrorDetail) {
	if bailServiceObj == nil {

		bailRepo, err := persistance.InitBailPersistance()
		if err != nil {
			return nil, err
		}

		bailServiceObj = &BailService{
			bailRepo: bailRepo,
		}
	}
	return bailServiceObj, nil
}
func (svc *BailService) GetBailInfoByQuality(request commonModels.BailInfoReuest) commonModels.BailInfoResponse {
	purchaseDetailes, purchaseErr := svc.bailRepo.GetPurchasedBailDetailByQuanlity(request.BranchId, request.Quality)
	if purchaseErr != nil {
		return commonModels.BailInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: "Error in getting purchase detailes for mentioned bail number",
				Errors: []commonModels.ErrorDetail{
					*purchaseErr,
				},
			},
		}
	}
	return commonModels.BailInfoResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Purchase: purchaseDetailes,
	}
}
func (svc *BailService) GetBailInfo(request commonModels.BailInfoReuest) commonModels.BailInfoResponse {
	purchaseDetailes, purchaseErr := svc.bailRepo.GetPurchasedBailDetail(request.BranchId, request.BailNo)
	if purchaseErr != nil {
		return commonModels.BailInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: "Error in getting bail info for mentioned bail number",
				Errors: []commonModels.ErrorDetail{
					*purchaseErr,
				},
			},
		}
	}

	salesDetailes, salesErr := svc.bailRepo.GetSalesBailDetail(request.BranchId, request.BailNo)
	if salesErr != nil && salesErr.ErrorCode != commonModels.ErrorNoDataFound {
		return commonModels.BailInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: "Error in getting sales detailes for mentioned bail number",
				Errors: []commonModels.ErrorDetail{
					*salesErr,
				},
			},
		}
	}
	return commonModels.BailInfoResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Purchase: []commonModels.BailDetailsDto{
			*purchaseDetailes,
		},
		Sales: salesDetailes,
	}
}

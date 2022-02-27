package service

import (
	commonModels "commonpkg/models"
	"item-service/persistance"
	"net/http"
)

var baleServiceObj *BaleService

type BaleService struct {
	baleRepo *persistance.BalePersistance
}

func InitBaleService() (*BaleService, *commonModels.ErrorDetail) {
	if baleServiceObj == nil {

		baleRepo, err := persistance.InitBalePersistance()
		if err != nil {
			return nil, err
		}

		baleServiceObj = &BaleService{
			baleRepo: baleRepo,
		}
	}
	return baleServiceObj, nil
}
func (svc *BaleService) GetBaleInfoByQuality(request commonModels.BaleInfoReuest) commonModels.BaleInfoResponse {
	purchaseDetailes, purchaseErr := svc.baleRepo.GetPurchasedBaleDetailByQuanlity(request.GodownId, request.Quality)
	if purchaseErr != nil {
		errMessage := "Error in getting purchase detailes for mentioned bale number"
		if purchaseErr.ErrorCode == commonModels.ErrorNoDataFound {
			errMessage = "No Purchase for mentioned Quality"
		}
		return commonModels.BaleInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: errMessage,
				Errors: []commonModels.ErrorDetail{
					*purchaseErr,
				},
			},
		}
	}
	return commonModels.BaleInfoResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Purchase: purchaseDetailes,
	}
}
func (svc *BaleService) GetBaleInfo(request commonModels.BaleInfoReuest) commonModels.BaleInfoResponse {
	purchaseDetailes, purchaseErr := svc.baleRepo.GetPurchasedBaleDetail(request.GodownId, request.BaleNo)
	if purchaseErr != nil {
		return commonModels.BaleInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: "Error in getting bale info for mentioned bale number",
				Errors: []commonModels.ErrorDetail{
					*purchaseErr,
				},
			},
		}
	}

	salesDetailes, salesErr := svc.baleRepo.GetSalesBaleDetail(request.GodownId, request.BaleNo, "")
	if salesErr != nil && salesErr.ErrorCode != commonModels.ErrorNoDataFound {
		return commonModels.BaleInfoResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: "Error in getting sales detailes for mentioned bale number",
				Errors: []commonModels.ErrorDetail{
					*salesErr,
				},
			},
		}
	}
	return commonModels.BaleInfoResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Purchase: []commonModels.BaleDetailsDto{
			*purchaseDetailes,
		},
		Sales: salesDetailes,
	}
}

package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"
	"item-service/persistance"
	"net/http"
	"strings"
)

var purchaseServiceObj *PurchaseService

type PurchaseService struct {
	purchaseRepo *persistance.PurchasePersistance
	baleRepo     *persistance.BalePersistance
}

func InitPurchaseService() (*PurchaseService, *commonModels.ErrorDetail) {
	if purchaseServiceObj == nil {
		purchaseRepo, err := persistance.InitPurchasePersistance()
		if err != nil {
			return nil, err
		}

		baleRepo, err := persistance.InitBalePersistance()
		if err != nil {
			return nil, err
		}

		purchaseServiceObj = &PurchaseService{
			purchaseRepo: purchaseRepo,
			baleRepo:     baleRepo,
		}
	}
	return purchaseServiceObj, nil
}

func (svc *PurchaseService) GetAllPurchaseOrders(request commonModels.InventoryListRequest) commonModels.InventoryListResponse {
	list, lastEvalutionKey, err := svc.purchaseRepo.GetAllPurchaseOrders(request)
	if err != nil {
		return commonModels.InventoryListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting Purchase orders",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	request.LastEvalutionKey = nil
	total, err := svc.purchaseRepo.GetTotalPurchaseOrders(request)
	if err != nil {
		return commonModels.InventoryListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting Purchase orders",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	return commonModels.InventoryListResponse{
		CommonListResponse: commonModels.CommonListResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},

			LastEvalutionKey: lastEvalutionKey,
			PageSize:         request.PageSize,
			Total:            total,
		},
		Data: list,
	}
}

func (svc *PurchaseService) GetPurchaseBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {

	data, err := svc.purchaseRepo.GetPurchaseBillDetails(request)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get details for purchase bill no %s", request.PurchaseBillNumber),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	return commonModels.InventoryResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: *data,
	}
}

func (svc *PurchaseService) AddPurchaseBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
	err := validPurchaseUpsertrequest(request)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add details for purchase bill no %s", request.BillNo),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	return upsertPurchaseBill(request, true)
}

func upsertPurchaseBill(request commonModels.InventoryDto, isAdd bool) commonModels.InventoryResponse {
	request.InventorySortKey = common.GetInventoryPurchanseSortKey(request.BillNo)
	_, err := purchaseServiceObj.purchaseRepo.UpsertPurchaseOrder(request)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	for _, val := range request.BaleDetails {
		islongation := false
		val.GodownId = request.GodownId
		val.BillNo = request.BillNo
		val.PurchaseDate = request.PurchaseDate
		if isAdd {
			val.PendingQuantity = val.BilledQuantity
		}
		val.SortKey = common.GetBaleDetailPurchanseSortKey(val.Quality, val.BaleNo)
		if val.ReceivedQuantity > 0 && val.ReceivedQuantity-val.BilledQuantity > 0 {
			islongation = true
		}

		baleInfo := commonModels.BaleInfoDto{
			GodownId:         request.GodownId,
			BaleInfoSortKey:  common.GetBaleInfoSortKey(val.BaleNo),
			BaleNo:           val.BaleNo,
			ReceivedQuantity: val.ReceivedQuantity,
			BilledQuantity:   val.BilledQuantity,
			IsLongation:      islongation,
			Quality:          val.Quality,
		}
		_, err := purchaseServiceObj.baleRepo.UpsertBaleInfo(baleInfo)
		if err != nil {
			return commonModels.InventoryResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			}
		}
		_, err = purchaseServiceObj.baleRepo.UpsertBaleDetail(val)
		if err != nil {
			return commonModels.InventoryResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			}
		}
	}
	return commonModels.InventoryResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusCreated,
		},
		Data: request,
	}
}
func validPurchaseUpsertrequest(request commonModels.InventoryDto) *commonModels.ErrorDetail {
	oldPurchaseBill, err := purchaseServiceObj.purchaseRepo.GetPurchaseBillDetails(commonModels.InventoryFilterDto{
		GodownId:           request.GodownId,
		PurchaseBillNumber: request.BillNo,
	})

	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
		}
	}
	if oldPurchaseBill != nil {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorAlreadyExists,
			ErrorMessage: fmt.Sprintf("same purchase bill no already exists, bill no %s", request.BillNo),
		}
	}

	errlist := make([]string, 0)
	for _, val := range request.BaleDetails {
		oldBaleInfo, err := purchaseServiceObj.baleRepo.GetPurchasedBaleDetail(request.GodownId, val.BaleNo)

		if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
			return &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorServer,
				ErrorMessage: fmt.Sprintf("could not add details for purchase bill no %s", request.BillNo),
			}
		}
		if oldBaleInfo != nil {
			errlist = append(errlist, val.BaleNo)
		}

	}

	if len(errlist) > 0 {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not add details for purchase bill/bills [%s]", strings.Join(errlist, ", ")),
		}
	}
	return nil
}

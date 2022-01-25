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
	bailRepo     *persistance.BailPersistance
}

func InitPurchaseService() (*PurchaseService, *commonModels.ErrorDetail) {
	if purchaseServiceObj == nil {
		purchaseRepo, err := persistance.InitPurchasePersistance()
		if err != nil {
			return nil, err
		}

		bailRepo, err := persistance.InitBailPersistance()
		if err != nil {
			return nil, err
		}

		purchaseServiceObj = &PurchaseService{
			purchaseRepo: purchaseRepo,
			bailRepo:     bailRepo,
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

func (svc *PurchaseService) UpdatePurchaseBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
	err, bailstoBeDeleted := validPurchaseUpsertrequest(request, false)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not update details for purchase bill no %s", request.BillNo),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	for _, val := range bailstoBeDeleted {
		svc.bailRepo.DeleteBailDetails(request.BranchId, val)
		svc.bailRepo.DeleteBailInfo(request.BranchId, val)

	}
	return upsertPurchaseBill(request, false)
}

func (svc *PurchaseService) AddPurchaseBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
	err, _ := validPurchaseUpsertrequest(request, true)
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

func (svc *PurchaseService) DeletePurchaseBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {
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

	for _, val := range data.BailDetails {
		svc.bailRepo.DeleteBailDetails(request.BranchId, val.BailNo)
		svc.bailRepo.DeleteBailInfo(request.BranchId, val.BailNo)
	}
	svc.purchaseRepo.DeletePurchaseBillDetails(request.BranchId, request.PurchaseBillNumber)

	return commonModels.InventoryResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: *data,
	}
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

	for _, val := range request.BailDetails {
		islongation := false
		val.BranchId = request.BranchId
		val.BillNo = request.BillNo
		val.PurchaseDate = request.PurchaseDate
		if isAdd {
			val.PendingQuantity = val.BilledQuantity
		}
		val.SortKey = common.GetBailDetailPurchanseSortKey(val.Quality, val.BailNo)
		if val.ReceivedQuantity > 0 && val.ReceivedQuantity-val.BilledQuantity > 0 {
			islongation = true
		}

		bailInfo := commonModels.BailInfoDto{
			BranchId:         request.BranchId,
			BailInfoSortKey:  common.GetBailInfoSortKey(val.BailNo),
			BailNo:           val.BailNo,
			ReceivedQuantity: val.ReceivedQuantity,
			BilledQuantity:   val.BilledQuantity,
			IsLongation:      islongation,
			Quality:          val.Quality,
		}
		_, err := purchaseServiceObj.bailRepo.UpsertBailInfo(bailInfo)
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
		_, err = purchaseServiceObj.bailRepo.UpsertBailDetail(val)
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
func validPurchaseUpsertrequest(request commonModels.InventoryDto, isNew bool) (*commonModels.ErrorDetail, []string) {
	oldPurchaseBill, err := purchaseServiceObj.purchaseRepo.GetPurchaseBillDetails(commonModels.InventoryFilterDto{
		BranchId:           request.BranchId,
		PurchaseBillNumber: request.BillNo,
	})

	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
		}, nil
	}
	if isNew && oldPurchaseBill != nil {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorAlreadyExists,
			ErrorMessage: fmt.Sprintf("same purchase bill no already exists, bill no %s", request.BillNo),
		}, nil
	}

	errlist := make([]string, 0)
	for _, val := range request.BailDetails {
		oldBailInfo, err := purchaseServiceObj.bailRepo.GetPurchasedBailDetail(request.BranchId, val.BailNo)

		if isNew && err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
			return &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorServer,
				ErrorMessage: fmt.Sprintf("could not add details for purchase bill no %s", request.BillNo),
			}, nil
		}
		if (isNew && oldBailInfo != nil) || (!isNew && oldBailInfo != nil && oldBailInfo.BillNo != request.BillNo) {
			errlist = append(errlist, val.BailNo)
		}

	}

	if len(errlist) > 0 {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not add details for purchase bill/bills [%s]", strings.Join(errlist, ", ")),
		}, nil
	}
	var tobeDeleted []string
	if !isNew {
		tobeDeleted = getBailsToBeDeleted(*oldPurchaseBill, request)
	}
	return nil, tobeDeleted
}

func getBailsToBeDeleted(oldPurchaseBill, newPurchaseBill commonModels.InventoryDto) []string {
	var deleteBailNo = make([]string, 0)

	for _, oldBail := range oldPurchaseBill.BailDetails {
		found := false
		for _, newBail := range newPurchaseBill.BailDetails {
			if newBail.BailNo == oldBail.BailNo {
				found = true
			}
		}
		if !found {
			deleteBailNo = append(deleteBailNo, oldBail.BailNo)
		}
	}
	return deleteBailNo
}

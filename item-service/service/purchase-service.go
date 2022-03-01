package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"
	"item-service/persistance"
	"net/http"

	uuid "github.com/iris-contrib/go.uuid"
)

type IPurchaseService interface {
	GetAll(request commonModels.InventoryListRequest) commonModels.PurchaseListResponse
	GetbyId(purchaseId string) commonModels.PurchaseResponse
	Add(data commonModels.AddPurchaseDataRequest) commonModels.AddPurchaseDataResponse
}

var purchaseServiceObj *PurchaseService

type PurchaseService struct {
	purchaseRepo persistance.IPurchasePersistance
	baleRepo     persistance.IBalePersistance
}

func InitPurchaseService() (IPurchaseService, *commonModels.ErrorDetail) {
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
func (svc *PurchaseService) GetAll(request commonModels.InventoryListRequest) commonModels.PurchaseListResponse {
	list, lastEvalutionKey, err := svc.purchaseRepo.GetAll(request)
	if err != nil {
		return commonModels.PurchaseListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting basic challans",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	request.LastEvalutionKey = nil
	total, err := svc.purchaseRepo.GetAllTotal(request)
	if err != nil {
		return commonModels.PurchaseListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting basic challans",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	return commonModels.PurchaseListResponse{
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

func (svc *PurchaseService) GetbyId(purchaseId string) commonModels.PurchaseResponse {
	data, err := svc.purchaseRepo.GetById(purchaseId)
	if err != nil {
		return commonModels.PurchaseResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get details for basic challan id %s", purchaseId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	return commonModels.PurchaseResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: *data,
	}
}
func (svc *PurchaseService) Add(data commonModels.AddPurchaseDataRequest) commonModels.AddPurchaseDataResponse {
	var errors []commonModels.ErrorDetail = make([]commonModels.ErrorDetail, 0)
	id, _ := uuid.NewV1()
	data.PurchaseDetails.PurchaseId = id.String()
	data.PurchaseDetails.SortKey = common.GetPurchaseSortKey(data.PurchaseDetails.ProductId, data.PurchaseDetails.QualityId, data.PurchaseDetails.PurchaseId)

	data.PurchaseDetails.Status = common.STATUS_PURCHASE_STOCK

	for i := range data.BaleDetails {
		oldBale, err := svc.baleRepo.GetBaleInfoByBaleNo(data.BaleDetails[i].BaleNo)
		if oldBale != nil {
			errors = append(errors, commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorAlreadyExists,
				ErrorMessage: fmt.Sprintf("same bale no already exists, %s", data.BaleDetails[i].BaleNo),
			})
			continue
		}
		if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
			errors = append(errors, *err)
			continue
		}
		data.BaleDetails[i].PurchaseDetails = commonModels.BalePurchaseDetails{
			PurchaseId: data.PurchaseDetails.PurchaseId,
		}
		data.BaleDetails[i].ProductId = data.PurchaseDetails.ProductId
		data.BaleDetails[i].QualityId = data.PurchaseDetails.QualityId
		data.BaleDetails[i].SortKey = common.GetInStockBaleSortKey(data.PurchaseDetails.ProductId, data.PurchaseDetails.QualityId, data.BaleDetails[i].BaleNo)
	}

	if len(errors) > 0 {
		return commonModels.AddPurchaseDataResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error in addign basic challan due to multiple errors.",
				Errors:       errors,
			},
			PurchaseDetails: data.PurchaseDetails,
			BaleDetails:     data.BaleDetails,
		}
	}
	oldpurchase, getPurchaseBillNo := svc.purchaseRepo.GetByBillNo(data.PurchaseDetails.PurchaseBillNo)
	if oldpurchase != nil {
		return commonModels.AddPurchaseDataResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error in addign basic challan due to multiple errors.",
				Errors: []commonModels.ErrorDetail{
					{
						ErrorCode:    commonModels.ErrorAlreadyExists,
						ErrorMessage: fmt.Sprintf("same basic challan no %s already exists.", data.PurchaseDetails.PurchaseBillNo),
					},
				},
			},
			PurchaseDetails: data.PurchaseDetails,
			BaleDetails:     data.BaleDetails,
		}
	}
	if getPurchaseBillNo != nil && getPurchaseBillNo.ErrorCode != commonModels.ErrorNoDataFound {
		return commonModels.AddPurchaseDataResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error in addign basic challan due to multiple errors.",
				Errors: []commonModels.ErrorDetail{
					*getPurchaseBillNo,
				},
			},
			PurchaseDetails: data.PurchaseDetails,
			BaleDetails:     data.BaleDetails,
		}
	}
	_, err := svc.purchaseRepo.Add(data.PurchaseDetails)
	if err != nil {
		return commonModels.AddPurchaseDataResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error in addign basic challan due to multiple errors.",
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
			PurchaseDetails: data.PurchaseDetails,
			BaleDetails:     data.BaleDetails,
		}
	}
	batchinsertErr := svc.baleRepo.BatchInsertBale(data.BaleDetails)
	if batchinsertErr != nil {
		return commonModels.AddPurchaseDataResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error in addign basic challan due to multiple errors.",
				Errors: []commonModels.ErrorDetail{
					*batchinsertErr,
				},
			},
			PurchaseDetails: data.PurchaseDetails,
			BaleDetails:     data.BaleDetails,
		}
	}

	return commonModels.AddPurchaseDataResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		PurchaseDetails: data.PurchaseDetails,
		BaleDetails:     data.BaleDetails,
	}
}

// func (svc *PurchaseService) GetPurchaseBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {

// 	data, err := svc.purchaseRepo.GetPurchaseBillDetails(request)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusBadRequest,
// 				ErrorMessage: fmt.Sprintf("could not get details for purchase bill no %s", request.PurchaseBillNumber),
// 				Errors: []commonModels.ErrorDetail{
// 					*err,
// 				},
// 			},
// 		}
// 	}

// 	return commonModels.InventoryResponse{
// 		CommonResponse: commonModels.CommonResponse{
// 			StatusCode: http.StatusOK,
// 		},
// 		Data: *data,
// 	}
// }

// func (svc *PurchaseService) AddPurchaseBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
// 	err := validPurchaseUpsertrequest(request)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusBadRequest,
// 				ErrorMessage: fmt.Sprintf("could not add details for purchase bill no %s", request.BillNo),
// 				Errors: []commonModels.ErrorDetail{
// 					*err,
// 				},
// 			},
// 		}
// 	}

// 	return upsertPurchaseBill(request, true)
// }

// func upsertPurchaseBill(request commonModels.InventoryDto, isAdd bool) commonModels.InventoryResponse {
// 	request.InventorySortKey = common.GetInventoryPurchanseSortKey(request.BillNo)
// 	_, err := purchaseServiceObj.purchaseRepo.UpsertPurchaseOrder(request)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusBadRequest,
// 				ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
// 				Errors: []commonModels.ErrorDetail{
// 					*err,
// 				},
// 			},
// 		}
// 	}

// 	for _, val := range request.BaleDetails {
// 		islongation := false
// 		val.GodownId = request.GodownId
// 		val.BillNo = request.BillNo
// 		val.PurchaseDate = request.PurchaseDate
// 		if isAdd {
// 			val.PendingQuantity = val.BilledQuantity
// 		}
// 		val.SortKey = common.GetBaleDetailPurchanseSortKey(val.Quality, val.BaleNo)
// 		if val.ReceivedQuantity > 0 && val.ReceivedQuantity-val.BilledQuantity > 0 {
// 			islongation = true
// 		}

// 		baleInfo := commonModels.BaleInfoDto{
// 			GodownId:         request.GodownId,
// 			BaleInfoSortKey:  common.GetBaleInfoSortKey(val.BaleNo),
// 			BaleNo:           val.BaleNo,
// 			ReceivedQuantity: val.ReceivedQuantity,
// 			BilledQuantity:   val.BilledQuantity,
// 			IsLongation:      islongation,
// 			Quality:          val.Quality,
// 		}
// 		_, err := purchaseServiceObj.baleRepo.UpsertBaleInfo(baleInfo)
// 		if err != nil {
// 			return commonModels.InventoryResponse{
// 				CommonResponse: commonModels.CommonResponse{
// 					StatusCode:   http.StatusBadRequest,
// 					ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
// 					Errors: []commonModels.ErrorDetail{
// 						*err,
// 					},
// 				},
// 			}
// 		}
// 		_, err = purchaseServiceObj.baleRepo.UpsertBaleDetail(val)
// 		if err != nil {
// 			return commonModels.InventoryResponse{
// 				CommonResponse: commonModels.CommonResponse{
// 					StatusCode:   http.StatusBadRequest,
// 					ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
// 					Errors: []commonModels.ErrorDetail{
// 						*err,
// 					},
// 				},
// 			}
// 		}
// 	}
// 	return commonModels.InventoryResponse{
// 		CommonResponse: commonModels.CommonResponse{
// 			StatusCode: http.StatusCreated,
// 		},
// 		Data: request,
// 	}
// }
// func validPurchaseUpsertrequest(request commonModels.InventoryDto) *commonModels.ErrorDetail {
// 	oldPurchaseBill, err := purchaseServiceObj.purchaseRepo.GetPurchaseBillDetails(commonModels.InventoryFilterDto{
// 		GodownId:           request.GodownId,
// 		PurchaseBillNumber: request.BillNo,
// 	})

// 	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("could not add/update details for purchase bill no %s", request.BillNo),
// 		}
// 	}
// 	if oldPurchaseBill != nil {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorAlreadyExists,
// 			ErrorMessage: fmt.Sprintf("same purchase bill no already exists, bill no %s", request.BillNo),
// 		}
// 	}

// 	errlist := make([]string, 0)
// 	for _, val := range request.BaleDetails {
// 		oldBaleInfo, err := purchaseServiceObj.baleRepo.GetPurchasedBaleDetail(request.GodownId, val.BaleNo)

// 		if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 			return &commonModels.ErrorDetail{
// 				ErrorCode:    commonModels.ErrorServer,
// 				ErrorMessage: fmt.Sprintf("could not add details for purchase bill no %s", request.BillNo),
// 			}
// 		}
// 		if oldBaleInfo != nil {
// 			errlist = append(errlist, val.BaleNo)
// 		}

// 	}

// 	if len(errlist) > 0 {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("could not add details for purchase bill/bills [%s]", strings.Join(errlist, ", ")),
// 		}
// 	}
// 	return nil
// }

package service

import (
	commonModels "commonpkg/models"
	"item-service/persistance"
)

var salesServiceObj *SalesService

type SalesService struct {
	baleRepo persistance.IBalePersistance
	// salesRepo *persistance.SalesPersistance
}

func InitSalesService() (*SalesService, *commonModels.ErrorDetail) {
	if salesServiceObj == nil {

		baleRepo, err := persistance.InitBalePersistance()
		if err != nil {
			return nil, err
		}

		// salesRepo, err := persistance.InitSalesPersistance()
		// if err != nil {
		// 	return nil, err
		// }

		salesServiceObj = &SalesService{
			baleRepo: baleRepo,
			// salesRepo: salesRepo,
		}
	}
	return salesServiceObj, nil
}

// //updated
// func (svc *SalesService) GetAllSalesOrders(request commonModels.InventoryListRequest) commonModels.InventoryListResponse {
// 	list, lastEvalutionKey, err := svc.salesRepo.GetAllSalesOrders(request)
// 	if err != nil {
// 		return commonModels.InventoryListResponse{
// 			CommonListResponse: commonModels.CommonListResponse{
// 				CommonResponse: commonModels.CommonResponse{
// 					StatusCode:   http.StatusBadRequest,
// 					ErrorMessage: "Error in getting Purchase orders",
// 					Errors: []commonModels.ErrorDetail{
// 						*err,
// 					},
// 				},
// 			},
// 		}
// 	}
// 	request.LastEvalutionKey = nil
// 	total, err := svc.salesRepo.GetTotalSalesOrders(request)
// 	if err != nil {
// 		return commonModels.InventoryListResponse{
// 			CommonListResponse: commonModels.CommonListResponse{
// 				CommonResponse: commonModels.CommonResponse{
// 					StatusCode:   http.StatusBadRequest,
// 					ErrorMessage: "Error in getting Sales orders",
// 					Errors: []commonModels.ErrorDetail{
// 						*err,
// 					},
// 				},
// 			},
// 		}
// 	}
// 	return commonModels.InventoryListResponse{
// 		CommonListResponse: commonModels.CommonListResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode: http.StatusOK,
// 			},
// 			LastEvalutionKey: lastEvalutionKey,
// 			PageSize:         request.PageSize,
// 			Total:            total,
// 		},
// 		Data: list,
// 	}
// }

// func (svc *SalesService) GetSalesBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {

// 	data, err := svc.salesRepo.GetSalesBillDetails(request)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusBadRequest,
// 				ErrorMessage: fmt.Sprintf("could not get details for sales bill no %s", request.SalesBillNumber),
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

// //updated

// func (svc *SalesService) AddSalesBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
// 	err := validSalesUpsertrequest(request)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusBadRequest,
// 				ErrorMessage: fmt.Sprintf("could not add details for sales bill no %s", request.BillNo),
// 				Errors: []commonModels.ErrorDetail{
// 					*err,
// 				},
// 			},
// 		}
// 	}
// 	return upsertSalesBill(request, true)
// }

// func updatePurchaseBaleRemainingQuantity(godownId, baleNo string, reduceQuantity int32, deleingBaleSalesBillNo string) *commonModels.ErrorDetail {
// 	purchaseBale, err := salesServiceObj.baleRepo.GetPurchasedBaleDetail(godownId, baleNo)
// 	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 		return err
// 	}
// 	if purchaseBale != nil {
// 		purchaseBale.PendingQuantity = purchaseBale.BilledQuantity - reduceQuantity

// 		if purchaseBale.PendingQuantity <= 0 {
// 			err = salesServiceObj.baleRepo.UpdateBaleDetailsOutofStock(godownId, baleNo)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			_, err = salesServiceObj.baleRepo.UpdateBaleDetailQuantity(*purchaseBale)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	} else {

// 		outofStockBale, err := salesServiceObj.baleRepo.GetOutofStockBaleDetail(godownId, baleNo)
// 		if err != nil {
// 			return err
// 		}

// 		salesDetailes, _ := salesServiceObj.baleRepo.GetSalesBaleDetail(godownId, baleNo, "")
// 		var total int32 = 0

// 		for _, val := range salesDetailes {
// 			if len(deleingBaleSalesBillNo) == 0 || (len(deleingBaleSalesBillNo) > 0 && val.BillNo != deleingBaleSalesBillNo) {
// 				total = total + val.BilledQuantity
// 			}
// 		}

// 		outofStockBale.PendingQuantity = outofStockBale.BilledQuantity - total

// 		err = salesServiceObj.baleRepo.RegenrateOutofStockBale(*outofStockBale)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
// func updateRemainingQuantityForBale(godownId string, baleDetails []commonModels.BaleDetailsDto) *commonModels.ErrorDetail {
// 	for _, val := range baleDetails {
// 		salesbales, err := salesServiceObj.baleRepo.GetSalesBaleDetail(godownId, val.BaleNo, "")
// 		if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 			return err
// 		}

// 		var totalSales int32 = 0
// 		for _, salesBale := range salesbales {
// 			totalSales = totalSales + salesBale.BilledQuantity
// 		}
// 		err = updatePurchaseBaleRemainingQuantity(godownId, val.BaleNo, totalSales, "")
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func updateRemainingQuantity(godownId, salesBillNo string) *commonModels.ErrorDetail {
// 	data, err := salesServiceObj.salesRepo.GetSalesBillDetails(commonModels.InventoryFilterDto{
// 		GodownId:        godownId,
// 		SalesBillNumber: salesBillNo,
// 	})

// 	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 		return err
// 	}
// 	if err != nil && err.ErrorCode == commonModels.ErrorNoDataFound {
// 		data, err = salesServiceObj.salesRepo.GetDeletedSalesBillDetails(commonModels.InventoryFilterDto{
// 			GodownId:        godownId,
// 			SalesBillNumber: salesBillNo,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return updateRemainingQuantityForBale(godownId, data.BaleDetails)
// }
// func upsertSalesBill(request commonModels.InventoryDto, isAdd bool) commonModels.InventoryResponse {
// 	request.InventorySortKey = common.GetInventorySalesSortKey(request.BillNo)
// 	_, err := salesServiceObj.salesRepo.UpsertSalesOrder(request)
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

// 		val.GodownId = request.GodownId
// 		val.BillNo = request.BillNo
// 		val.SalesDate = request.SalesDate

// 		val.SortKey = common.GetBaleDetailSalesSortKey(val.Quality, val.BaleNo, request.BillNo)
// 		_, err = salesServiceObj.baleRepo.UpsertBaleDetail(val)
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
// 	err = updateRemainingQuantity(request.GodownId, request.BillNo)
// 	if err != nil {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode:   http.StatusConflict,
// 				ErrorMessage: "Error in updating the remaining quantity",
// 				Errors: []commonModels.ErrorDetail{
// 					*err,
// 				},
// 			},
// 			Data: request,
// 		}
// 	} else {
// 		return commonModels.InventoryResponse{
// 			CommonResponse: commonModels.CommonResponse{
// 				StatusCode: http.StatusOK,
// 			},
// 			Data: request,
// 		}
// 	}
// }
// func validSalesUpsertrequest(request commonModels.InventoryDto) *commonModels.ErrorDetail {
// 	oldSalesBill, err := salesServiceObj.salesRepo.GetSalesBillDetails(commonModels.InventoryFilterDto{
// 		GodownId:        request.GodownId,
// 		SalesBillNumber: request.BillNo,
// 	})

// 	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("could not add/update details for sales bill no %s", request.BillNo),
// 		}
// 	}
// 	if oldSalesBill != nil {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorAlreadyExists,
// 			ErrorMessage: fmt.Sprintf("same sales bill no already exists, bill no %s", request.BillNo),
// 		}
// 	}

// 	errlist := make([]string, 0)
// 	for _, val := range request.BaleDetails {
// 		oldBaleInfo, err := salesServiceObj.baleRepo.GetPurchasedBaleDetail(request.GodownId, val.BaleNo)
// 		if err != nil {
// 			return &commonModels.ErrorDetail{
// 				ErrorCode:    commonModels.ErrorServer,
// 				ErrorMessage: fmt.Sprintf("could not add details for sales bill no %s, err: %s", request.BillNo, err.Error()),
// 			}
// 		}

// 		if oldBaleInfo.PendingQuantity < val.BilledQuantity {
// 			errlist = append(errlist, val.BaleNo)
// 		}

// 	}

// 	if len(errlist) > 0 {
// 		return &commonModels.ErrorDetail{
// 			ErrorCode:    commonModels.ErrorServer,
// 			ErrorMessage: fmt.Sprintf("could not sale more that what we have in the balel for bill/bills [%s]", strings.Join(errlist, ", ")),
// 		}
// 	}

// 	return nil
// }

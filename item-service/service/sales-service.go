package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"item-service/common"
	"item-service/persistance"
	"net/http"
	"strings"
)

var salesServiceObj *SalesService

type SalesService struct {
	bailRepo  *persistance.BailPersistance
	salesRepo *persistance.SalesPersistance
}

func InitSalesService() (*SalesService, *commonModels.ErrorDetail) {
	if salesServiceObj == nil {

		bailRepo, err := persistance.InitBailPersistance()
		if err != nil {
			return nil, err
		}

		salesRepo, err := persistance.InitSalesPersistance()
		if err != nil {
			return nil, err
		}

		salesServiceObj = &SalesService{
			bailRepo:  bailRepo,
			salesRepo: salesRepo,
		}
	}
	return salesServiceObj, nil
}

//updated
func (svc *SalesService) GetAllSalesOrders(request commonModels.InventoryListRequest) commonModels.InventoryListResponse {
	list, lastEvalutionKey, err := svc.salesRepo.GetAllSalesOrders(request)
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
	total, err := svc.salesRepo.GetTotalSalesOrders(request)
	if err != nil {
		return commonModels.InventoryListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Error in getting Sales orders",
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

func (svc *SalesService) GetSalesBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {

	data, err := svc.salesRepo.GetSalesBillDetails(request)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get details for sales bill no %s", request.SalesBillNumber),
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

//updated

func (svc *SalesService) AddSalesBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
	err, _ := validSalesUpsertrequest(request, true)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add details for sales bill no %s", request.BillNo),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	return upsertSalesBill(request, true)
}

func (svc *SalesService) UpdateSalesBillDetails(request commonModels.InventoryDto) commonModels.InventoryResponse {
	err, bailstoBeDeleted := validSalesUpsertrequest(request, false)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not update details for sales bill no %s", request.BillNo),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	for _, val := range bailstoBeDeleted {
		svc.bailRepo.DeleteSalesBailDetails(request.BranchId, val, request.BillNo)
	}
	return upsertSalesBill(request, false)
}

func (svc *SalesService) DeleteSalesBillDetails(request commonModels.InventoryFilterDto) commonModels.InventoryResponse {
	data, err := svc.salesRepo.GetSalesBillDetails(request)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get details for sales bill no %s", request.SalesBillNumber),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	for _, val := range data.BailDetails {
		svc.bailRepo.DeleteSalesBailDetails(request.BranchId, val.BailNo, request.SalesBillNumber)
	}
	svc.salesRepo.DeleteSalesBillDetails(request.BranchId, request.SalesBillNumber)
	err = updateRemainingQuantity(request.BranchId, request.SalesBillNumber)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusConflict,
				ErrorMessage: "Error in updating the remaining quantity",
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
			Data: *data,
		}
	} else {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: *data,
		}
	}

}
func updateRemainingQuantity(branchId, salesBillNo string) *commonModels.ErrorDetail {
	data, err := salesServiceObj.salesRepo.GetSalesBillDetails(commonModels.InventoryFilterDto{
		BranchId:        branchId,
		SalesBillNumber: salesBillNo,
	})

	if err != nil {
		return err
	}

	for _, val := range data.BailDetails {
		salesbails, err := salesServiceObj.bailRepo.GetSalesBailDetail(branchId, val.BailNo, salesBillNo)
		fmt.Printf("sales bails %+v \n", salesbails)
		if err != nil {
			return err
		}
		var totalSales int32 = 0
		for _, salesBail := range salesbails {
			totalSales = totalSales + salesBail.BilledQuantity
		}
		purchaseBail, err := salesServiceObj.bailRepo.GetPurchasedBailDetail(branchId, val.BailNo)
		if err != nil {
			return err
		}
		purchaseBail.PendingQuantity = purchaseBail.BilledQuantity - totalSales
		if purchaseBail.PendingQuantity <= 0 {
			purchaseBail.SortKey = common.GetBailDetailOutOfStockSortKey(purchaseBail.Quality, purchaseBail.BailNo)
			// err = salesServiceObj.bailRepo.DeleteBailDetails(branchId, val.BailNo)
			fmt.Printf("out of stock purchaseBail %+v,\n", purchaseBail)
			// if err != nil {
			// 	return err
			// }
		}
		fmt.Printf("purchaseBail %+v", purchaseBail)
		err = salesServiceObj.bailRepo.DeleteBailDetails(branchId, val.BailNo)
		if err != nil {
			return err
		}

		_, err = salesServiceObj.bailRepo.UpsertBailDetail(*purchaseBail)
		if err != nil {
			return err
		}
	}
	return nil
}
func upsertSalesBill(request commonModels.InventoryDto, isAdd bool) commonModels.InventoryResponse {
	request.InventorySortKey = common.GetInventorySalesSortKey(request.BillNo)
	_, err := salesServiceObj.salesRepo.UpsertSalesOrder(request)
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

		val.BranchId = request.BranchId
		val.BillNo = request.BillNo
		val.SalesDate = request.SalesDate

		val.SortKey = common.GetBailDetailSalesSortKey(val.Quality, val.BailNo, request.BillNo)
		_, err = salesServiceObj.bailRepo.UpsertBailDetail(val)
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
	err = updateRemainingQuantity(request.BranchId, request.BillNo)
	if err != nil {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusConflict,
				ErrorMessage: "Error in updating the remaining quantity",
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
			Data: request,
		}
	} else {
		return commonModels.InventoryResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			Data: request,
		}
	}
}
func validSalesUpsertrequest(request commonModels.InventoryDto, isNew bool) (*commonModels.ErrorDetail, []string) {
	oldSalesBill, err := salesServiceObj.salesRepo.GetSalesBillDetails(commonModels.InventoryFilterDto{
		BranchId:           request.BranchId,
		PurchaseBillNumber: request.BillNo,
	})

	if err != nil && err.ErrorCode != commonModels.ErrorNoDataFound {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not add/update details for sales bill no %s", request.BillNo),
		}, nil
	}
	if isNew && oldSalesBill != nil {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorAlreadyExists,
			ErrorMessage: fmt.Sprintf("same sales bill no already exists, bill no %s", request.BillNo),
		}, nil
	}

	errlist := make([]string, 0)
	for _, val := range request.BailDetails {
		oldBailInfo, err := salesServiceObj.bailRepo.GetPurchasedBailDetail(request.BranchId, val.BailNo)
		if err != nil {
			return &commonModels.ErrorDetail{
				ErrorCode:    commonModels.ErrorServer,
				ErrorMessage: fmt.Sprintf("could not add details for sales bill no %s, err: %s", request.BillNo, err.Error()),
			}, nil
		}
		if oldBailInfo.PendingQuantity < val.BilledQuantity {
			errlist = append(errlist, val.BailNo)
		}

	}

	if len(errlist) > 0 {
		return &commonModels.ErrorDetail{
			ErrorCode:    commonModels.ErrorServer,
			ErrorMessage: fmt.Sprintf("could not sale more that what we have in the baill for bill/bills [%s]", strings.Join(errlist, ", ")),
		}, nil
	}
	var tobeDeleted []string
	if !isNew {
		tobeDeleted = getBailsToBeDeleted(*oldSalesBill, request)
	}
	return nil, tobeDeleted
}

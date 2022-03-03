import { CommonListResponse, CommonResponse } from "./genric-model"


export interface AddPurchaseDataRequest {
    purchaseDetails: PurchaseDto;
    baleDetails: BaleDetailsDto[];
}
export interface AddPurchaseDataResponse extends CommonResponse {
    purchaseDetails: PurchaseDto;
    baleDetails: BaleDetailsDto[];
}

export interface PurchaseDto {
    godownId: string;
    sortKey: string; /// ProductId|QualityId|PurchaseId
    purchaseId: string; // GSI -1  PK  (all attr)
    purchaseBillNo: string; // GSI - 2 PK (keys only)
    date: Date;
    vendorId: string;
    vendorName: string;
    productId: string;
    productName: string;
    qualityId: string;
    qualityName: string;
    hsnCode: string;
    purchaseStatus: string; // Stock | Sold
}

export interface SalesDto {
    godownId: string;
    sortKey: string;// ProductId|QualityId|SalesId
    salesId: string;// GSI -1  PK  (all attr)
    salesBillNo: string;// GSI PK (keys only)
    clientId: string;
    transporterId: string;
    lrNo: string;
    challanNo: string; // GSI PK (keys only)
    date: Date;
    productId: string;
    qualityId: string;
}

export interface BalePurchaseDetails {
    purchaseId: string;
}
export interface BaleSalesDetails {
    salesId: string;
}
export interface BaleTransferDetails {
    fromGodownId: string;
    toGowodnId: string;
    date: Date;
}

export interface BaleDetailsDto {
    godownId: string;
    sortKey: string;
    baleNo: string;
    productId: string;
    qualityId: string;
    billedQuantity: number;
    receivedQuantity: number;
    rate: number;
    purchaseDetails: BalePurchaseDetails;
    salesDetails: BaleSalesDetails;
    transferDetails: BaleTransferDetails[];
}

export interface PurchaseListResponse extends CommonListResponse{
	data:PurchaseDto[];
}

export interface BailDetailsDto {
    branchId: string;
    sortKey: string; //// Bail | <Purchase or Sales or OutOfStock>|Quality|BailNo|<salesBill or purchseBill number>
    bailNo: string;
    quality: string;
    isSales: boolean;
    billNo: string;
    rate: number;
    purchaseDate: Date;
    salesDate: Date;
    clientId: string;
    vendorId: string;
    transferedToBranchId: string;
    transferedFromBranchId: string;
    receivedQuantity: number;
    billedQuantity: number;
    pendingQuantity: number;
}

export interface BailInfoDto {
    branchId: string;
    bailInfoSortKey: string; /// Info | bailNo | Quality
    bailNo: string;
    receivedQuantity: number;
    billedQuantity: number;
    isLongation: boolean;
    quality: string;
}

export interface InventoryDto {
    branchId: string;
    inventorySortKey: string; /// Inventory | <Purchase or Sales>| Bill No
    billNo: string;
    bailDetails: BailDetailsDto[];
    purchaseDate: Date;
    salesDate: Date;
    transporterId: string;
    lrNo: string;
    challanNo: string;
    hsnCode: string;
}

export interface InventoryFilterDto {
    godownId: string
    purchaseBillNumber: string;
    salesBillNumber: string;
    qualityId: string;
    productId: string;
}

export interface InventoryListRequest extends InventoryFilterDto {
    lastEvalutionKey: any;
    pageSize: number;

}
export interface InventoryListResponse extends CommonListResponse {
    data: InventoryDto[];
}

export interface InventoryResponse extends CommonResponse {
    data: InventoryDto;
}

export interface BailInfoReuest {
    branchId: string;
    bailNo: string;
    quality: string;
}

export interface BailInfoResponse extends CommonResponse {
    purchase: BailDetailsDto[];
    sales: BailDetailsDto[];
}

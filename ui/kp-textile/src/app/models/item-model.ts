import { CommonListResponse, CommonResponse } from "./genric-model"

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
    branchId: string
    purchaseBillNumber: string;
    salesBillNumber: string;
    startDate: Date;
    endDate: Date;
    quality: string;
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

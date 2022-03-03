import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddPurchaseDataRequest, AddPurchaseDataResponse, InventoryDto, InventoryListRequest, InventoryListResponse, InventoryResponse, PurchaseListResponse } from "../models/item-model";
import { environment } from "src/environments/environment";


@Injectable({
    providedIn: 'root',
})
export class PurchaseService {
    baseUrl: string = environment.purchaseBaseUrl
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllPurchase(lastEvalutionKey: any, pageSize: number, godownId: string):Observable<PurchaseListResponse>{
        return this.httpClient.post<PurchaseListResponse>(`${this.baseUrl}/getall`, { 
            lastEvalutionKey: lastEvalutionKey,
            pageSize: pageSize,
            godownId: godownId
         } as InventoryListRequest);
    }
    addPurchase(data: AddPurchaseDataRequest):Observable<AddPurchaseDataResponse>{
        return this.httpClient.post<AddPurchaseDataResponse>(`${this.baseUrl}/`, data);
    }

}
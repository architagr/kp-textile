import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { InventoryDto, InventoryListRequest, InventoryListResponse, InventoryResponse } from "../models/item-model";


@Injectable({
    providedIn: 'root',
})
export class PurchaseService {
    baseUrl: string = "http://localhost:8084/purchase"
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllPurchase(lastEvalutionKey: any, pageSize: number):Observable<InventoryListResponse>{
        return this.httpClient.post<InventoryListResponse>(`${this.baseUrl}/getall`, { 
            lastEvalutionKey: lastEvalutionKey,
            pageSize: pageSize,
         } as InventoryListRequest);
    }
    addPurchase(data: InventoryDto):Observable<InventoryResponse>{
        return this.httpClient.post<InventoryResponse>(`${this.baseUrl}/`, data);
    }
    updatePurchase(data: InventoryDto):Observable<InventoryResponse>{
        return this.httpClient.put<InventoryResponse>(`${this.baseUrl}/${data.billNo}`, data);
    }
    deletePurchaseOrder(billNumber: string):Observable<InventoryResponse>{
        return this.httpClient.delete<InventoryResponse>(`${this.baseUrl}/${billNumber}`);
    }
    getPurchaseOrder(billNumber: string):Observable<InventoryResponse>{
        return this.httpClient.get<InventoryResponse>(`${this.baseUrl}/${billNumber}`);
    }

}
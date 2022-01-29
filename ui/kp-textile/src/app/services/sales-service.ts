import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { InventoryDto, InventoryListRequest, InventoryListResponse, InventoryResponse } from "../models/item-model";
import { environment } from "src/environments/environment";


@Injectable({
    providedIn: 'root',
})
export class SalesService {
    baseUrl: string = environment.salesBaseUrl;
    constructor(
        private httpClient: HttpClient
    ) { }

    getAllSales(lastEvalutionKey: any, pageSize: number):Observable<InventoryListResponse>{
        return this.httpClient.post<InventoryListResponse>(`${this.baseUrl}/getall`, { 
            lastEvalutionKey: lastEvalutionKey,
            pageSize: pageSize,
         } as InventoryListRequest);
    }
    addSales(data: InventoryDto):Observable<InventoryResponse>{
        return this.httpClient.post<InventoryResponse>(`${this.baseUrl}/`, data);
    }
    updateSales(data: InventoryDto):Observable<InventoryResponse>{
        return this.httpClient.put<InventoryResponse>(`${this.baseUrl}/${data.billNo}`, data);
    }
    deleteSales(billNumber: string):Observable<InventoryResponse>{
        return this.httpClient.delete<InventoryResponse>(`${this.baseUrl}/${billNumber}`);
    }
    getSales(billNumber: string):Observable<InventoryResponse>{
        return this.httpClient.get<InventoryResponse>(`${this.baseUrl}/${billNumber}`);
    }
}
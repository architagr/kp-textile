import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { InventoryDto, InventoryListRequest, InventoryListResponse, InventoryResponse } from "../models/item-model";


@Injectable({
    providedIn: 'root',
})
export class DocumentService {
    baseUrl: string = "http://localhost:8085/"
    constructor(
        private httpClient: HttpClient
    ) { }


    getChallan(data: InventoryDto):Observable<any> {
        return this.httpClient.post(`${this.baseUrl}challan`, data, {responseType: 'text'})
    }
}
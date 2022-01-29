import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { InventoryDto } from "../models/item-model";
import { environment } from "src/environments/environment";

@Injectable({
    providedIn: 'root',
})
export class DocumentService {
    baseUrl: string = environment.documentBaseUrl
    constructor(
        private httpClient: HttpClient
    ) { }


    getChallan(data: InventoryDto):Observable<any> {
        return this.httpClient.post(`${this.baseUrl}challan`, data, {responseType: 'text'})
    }
}
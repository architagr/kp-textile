import { Injectable } from "@angular/core";
import { Observable, of } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddRequest, HnsCodeDto, HnsCodeListResponse, HnsCodeResponse } from "../models/hsn-code-model";
import { ProductDto, ProductListResponse, ProductResponse, QualityDto, QualityListResponse, QualityResponse } from "../models/quality-model";
import { environment } from "src/environments/environment";
import { GodownListResponse, GodownResponse } from "../models/godown-model";

@Injectable({
    providedIn: 'root',
})
export class GodownService {
    baseUrl: string = `${environment.organizationBaseUrl}godown/`;

    constructor(
        private httpClient: HttpClient
    ) { }

    getAllGodown(): Observable<GodownListResponse> {
        return this.httpClient.get<GodownListResponse>(`${this.baseUrl}`)
    }
    addGodown(name: string): Observable<GodownResponse> {
        return this.httpClient.post<GodownResponse>(`${this.baseUrl}`, {name})
    }
}
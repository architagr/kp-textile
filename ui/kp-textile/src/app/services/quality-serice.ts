import { Injectable } from "@angular/core";
import { Observable, of } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddRequest, HnsCodeDto, HnsCodeListResponse, HnsCodeResponse } from "../models/hsn-code-model";
import { ProductDto, ProductListResponse, ProductResponse, QualityDto, QualityListResponse, QualityResponse } from "../models/quality-model";
import { environment } from "src/environments/environment";

@Injectable({
    providedIn: 'root',
})
export class QualityService {
    baseUrl: string = environment.qualityBaseUrl;
    constructor(
        private httpClient: HttpClient
    ) { }

    getAllQualities(): Observable<QualityListResponse> {
        return this.httpClient.get<QualityListResponse>(`${this.baseUrl}quality/`)
    }
    addQuality(quality: QualityDto): Observable<QualityResponse> {
        return this.httpClient.post<QualityResponse>(`${this.baseUrl}quality/`, quality)
    }
    addProduct(name: string): Observable<ProductResponse> {
        return this.httpClient.post<ProductResponse>(`${this.baseUrl}product/`, { code: name })
    }

    getAllProduct(): Observable<ProductListResponse> {
        return this.httpClient.get<ProductListResponse>(`${this.baseUrl}product/`)
    }
}
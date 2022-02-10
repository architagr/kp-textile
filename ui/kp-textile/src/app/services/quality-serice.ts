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

    productId: number = 5
    products: ProductDto[] = [
        { id: '1', name: "Rayon" },
        { id: '2', name: "Silk" },
        { id: '3', name: "Cotton" },
        { id: '4', name: "Woolen" },
    ]

    qualityId: number = 7
    qualities: QualityDto[] = [
        { id: '1', name: "01-01", hsnCode: "hsn01-01", productId: '1', productName: '' },
        { id: '2', name: "01-02", hsnCode: "hsn01-02", productId: '1', productName: '' },
        { id: '3', name: "03-03", hsnCode: "hsn03-03", productId: '3', productName: '' },
        { id: '4', name: "02-04", hsnCode: "hsn02-04", productId: '2', productName: '' },
        { id: '5', name: "02-05", hsnCode: "hsn02-05", productId: '2', productName: '' },
        { id: '6', name: "04-06", hsnCode: "hsn04-06", productId: '4', productName: '' },
    ]

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
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddRequest, HnsCodeDto, HnsCodeListResponse, HnsCodeResponse } from "../models/hsn-code-model";
import { QualityDto, QualityListResponse, QualityResponse } from "../models/quality-model";

@Injectable({
    providedIn: 'root',
})
export class QualityService {
    baseUrl: string = "http://localhost:8080/"
    constructor(
        private httpClient: HttpClient
    ) { }

    getAllQualities(): Observable<QualityListResponse> {
        return this.httpClient.get<QualityListResponse>(`${this.baseUrl}`);
    }
    addQuality(quality: QualityDto): Observable<QualityResponse> {
        let body = { code: quality.name } as AddRequest
        return this.httpClient.post<QualityResponse>(`${this.baseUrl}`, body);
    }
}
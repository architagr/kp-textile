import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddRequest, HnsCodeDto, HnsCodeListResponse, HnsCodeResponse } from "../models/hsn-code-model";
import { environment } from "src/environments/environment";

@Injectable({
    providedIn: 'root',
})
export class HsnCodeService {
    private baseUrl: string = environment.hsnBaseUrl;
    constructor(
        private httpClient: HttpClient
    ) { }

    getAllHsnCode(): Observable<HnsCodeListResponse> {
        return this.httpClient.get<HnsCodeListResponse>(`${this.baseUrl}`);
    }
    addHsnCode(hsnCode: HnsCodeDto): Observable<HnsCodeResponse> {
        let body = { code: hsnCode.hnsCode } as AddRequest
        return this.httpClient.post<HnsCodeResponse>(`${this.baseUrl}`, body);
    }
}
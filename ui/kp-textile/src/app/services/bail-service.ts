import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { BailInfoResponse } from "../models/item-model";
import { environment } from "src/environments/environment";


@Injectable({
    providedIn: 'root',
})
export class BailService {
    baseUrl: string = environment.bailBaseUrl
    constructor(
        private httpClient: HttpClient
    ) { }


    getBailInfoByQuality(qualityId: string): Observable<BailInfoResponse> {
        return this.httpClient.get<BailInfoResponse>(`${this.baseUrl}/quality/${qualityId}`)
    }

    getBailInfo(baleNumber: string): Observable<BailInfoResponse> {
        return this.httpClient.get<BailInfoResponse>(`${this.baseUrl}/${baleNumber}`)
    }
}
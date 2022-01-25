import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { BailInfoResponse } from "../models/item-model";


@Injectable({
    providedIn: 'root',
})
export class BailService {
    baseUrl: string = "http://localhost:8084/bailInfo"
    constructor(
        private httpClient: HttpClient
    ) { }


    getBailInfoByQuality(qualityId: string): Observable<BailInfoResponse> {
        return this.httpClient.get<BailInfoResponse>(`${this.baseUrl}/quality/${qualityId}`)
    }
}
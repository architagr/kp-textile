import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { CommonResponse } from "../models/genric-model";
import { AddTransporterRequest, AddTransporterResponse, TransporterListResponse } from "../models/transporter-model";
import { environment } from "src/environments/environment";

@Injectable({
    providedIn: 'root',
})
export class TransporterService {
    private baseUrl: string = `${environment.transporterBaseUrl}transporter/`;
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllTransporter(pageSize: number, searchText: string, lastEvalutionKey: any | null): Observable<TransporterListResponse> {
        let url = `${this.baseUrl}getall?pageSize=${pageSize}`
        if (searchText.length > 0) {
            url += `&companyName=${searchText}`
        }
        return this.httpClient.post<TransporterListResponse>(url, { lastEvalutionKey: lastEvalutionKey });
    }

    addTransporter(client: AddTransporterRequest): Observable<AddTransporterResponse> {
        return this.httpClient.post<AddTransporterResponse>(`${this.baseUrl}`, client);
    }

    getTransporterData(vendorId: string): Observable<AddTransporterResponse> {
        return this.httpClient.get<AddTransporterResponse>(`${this.baseUrl}${vendorId}`);
    }

    updateTransporter(vendorId: string, data: AddTransporterResponse): Observable<AddTransporterResponse> {
        return this.httpClient.put<AddTransporterResponse>(`${this.baseUrl}${vendorId}`, data);
    }
    deleteTransporter(vendorId: string): Observable<CommonResponse> {
        return this.httpClient.delete<CommonResponse>(`${this.baseUrl}${vendorId}`);
    }
}
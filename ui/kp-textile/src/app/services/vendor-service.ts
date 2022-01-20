import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { CommonResponse } from "../models/genric-model";
import { AddVendorRequest, AddVendorResponse, VendorListResponse } from "../models/vendor-models";

@Injectable({
    providedIn: 'root',
})
export class VendorService {
    baseUrl: string = "http://localhost:8082/"
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllVendors(pageSize: number, searchText: string, lastEvalutionKey: any | null): Observable<VendorListResponse> {
        let url = `${this.baseUrl}getall?pageSize=${pageSize}`
        if (searchText.length > 0) {
            url += `&companyName=${searchText}`
        }
        return this.httpClient.post<VendorListResponse>(url, { lastEvalutionKey: lastEvalutionKey });
    }

    addVendor(client: AddVendorRequest): Observable<AddVendorResponse> {
        return this.httpClient.post<AddVendorResponse>(`${this.baseUrl}`, client);
    }

    getVendorData(vendorId: string): Observable<AddVendorResponse> {
        return this.httpClient.get<AddVendorResponse>(`${this.baseUrl}${vendorId}`);
    }

    updateVendor(vendorId: string, data: AddVendorResponse): Observable<AddVendorResponse> {
        return this.httpClient.put<AddVendorResponse>(`${this.baseUrl}${vendorId}`, data);
    }
    deleteVendor(vendorId: string): Observable<CommonResponse> {
        return this.httpClient.delete<CommonResponse>(`${this.baseUrl}${vendorId}`);
    }
}
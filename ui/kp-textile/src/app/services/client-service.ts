import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddClientRequest, AddClientResponse, ClientListResponse } from "../models/client-model";
import { CommonResponse } from "../models/genric-model";

@Injectable({
    providedIn: 'root',
})
export class ClientService {
    baseUrl: string = "http://localhost:8080/"
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllClient(pageSize: number, searchText: string, lastEvalutionKey: any | null): Observable<ClientListResponse> {
        let url = `${this.baseUrl}getall?pageSize=${pageSize}`
        if (searchText.length > 0) {
            url += `&companyName=${searchText}`
        }
        return this.httpClient.post<ClientListResponse>(url, { lastEvalutionKey: lastEvalutionKey });
    }

    addClient(client: AddClientRequest): Observable<AddClientResponse> {
        return this.httpClient.post<AddClientResponse>(`${this.baseUrl}`, client);
    }

    getClientData(clientId: string): Observable<AddClientResponse> {
        return this.httpClient.get<AddClientResponse>(`${this.baseUrl}${clientId}`);
    }

    updateClient(clientId: string, data: AddClientResponse): Observable<AddClientResponse> {
        return this.httpClient.put<AddClientResponse>(`${this.baseUrl}${clientId}`, data);
    }
    deleteClient(clientId: string): Observable<CommonResponse> {
        return this.httpClient.delete<CommonResponse>(`${this.baseUrl}${clientId}`);
    }
}
import { Injectable } from "@angular/core";
import { map, Observable, Subject, switchMap } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { AddClientRequest, AddClientResponse, ClientDto, ClientListResponse, ClientResponse } from "../models/client-model";

@Injectable({
    providedIn: 'root',
})
export class ClientService {
    baseUrl:string = "http://localhost:8080/"
    constructor(
        private httpClient: HttpClient
    ) { }


    getAllClient(): Observable<ClientListResponse> {
        return this.httpClient.post<ClientListResponse>(`${this.baseUrl}getall?pageSize=10`, {});
    }

    addClient(client: AddClientRequest): Observable<AddClientResponse>{
        return this.httpClient.post<AddClientResponse>(`${this.baseUrl}`, client);
    }
}
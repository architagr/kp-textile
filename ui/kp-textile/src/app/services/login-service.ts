import { Injectable } from "@angular/core";
import { Observable, of } from "rxjs";
import { HttpClient } from "@angular/common/http"
import { environment } from "src/environments/environment";
import { LoginRequest, LoginResponse } from "../models/user-model";

@Injectable({
    providedIn: 'root',
})
export class LoginService {
    private baseUrl: string = `${environment.organizationBaseUrl}user/`;

    constructor(
        private httpClient: HttpClient
    ) { }

    login(loginRequest: LoginRequest): Observable<LoginResponse> {
        return this.httpClient.post<LoginResponse>(`${this.baseUrl}login`, loginRequest);
    }
 
}
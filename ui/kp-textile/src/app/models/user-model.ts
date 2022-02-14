import { CommonResponse } from "./genric-model";

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse extends CommonResponse {
    token: string;
}
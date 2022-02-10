import { Injectable } from "@angular/core";

@Injectable({
    providedIn: 'root',
})
export class SpinnerService {
    count: number = 0;
    constructor(
    ) { }

    show(){
        this.count++;
    }
    hide(){
        this.count--;
        if(this.count<0){
            this.count = 0;
        }
    }
}
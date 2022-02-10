import { Injectable } from "@angular/core";
import { ToastrService } from "ngx-toastr";

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(private toastr: ToastrService,
    ) { }

    show(type: 'Success' | 'Error', message: string) {
        let toastClass = 'alert-success'
        if (type === "Error") {
            toastClass = 'alert-danger'
        }
        this.toastr.info(`<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> ${message}.`, type, {
            disableTimeOut: false,
            timeOut: 3000,
            closeButton: true,
            enableHtml: true,
            toastClass: `alert ${toastClass} alert-with-icon`,
            positionClass: 'toast-top-right'
        });

    }
}
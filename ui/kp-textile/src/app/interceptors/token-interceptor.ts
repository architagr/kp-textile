import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpEvent, HttpResponse, HttpRequest, HttpHandler } from '@angular/common/http';
import { catchError, Observable, of, tap } from 'rxjs';
import { SpinnerService } from '../services/spinner-service';
import { ToastService } from '../services/toast-service';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(private spinnerService: SpinnerService,
  ) { }
  intercept(httpRequest: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    this.spinnerService.show();
    const Authorization = `Bearer ${localStorage.getItem('token')}`;
    return next.handle(httpRequest.clone({ setHeaders: { Authorization } })).pipe(
      tap(() => {
          this.spinnerService.hide();
      }),
      catchError((err: any) => {
        console.log(`error`, { err })
        
        return of(err);
      })
    );
  }
}
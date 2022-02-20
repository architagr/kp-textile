import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { Route, Router } from '@angular/router';
import { LoginRequest, LoginResponse } from 'src/app/models/user-model';
import { LoginService } from 'src/app/services/login-service';
import { ToastService } from 'src/app/services/toast-service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  loginForm: FormGroup
  constructor(private loginService: LoginService,
    private toasterService: ToastService,
    private router: Router) {
    this.loginForm = new FormGroup({
      username: new FormControl('', [Validators.required]),
      password: new FormControl('', [Validators.required]),
    })
  }
  
  get formControls(): { [key: string]: AbstractControl } {
    return this.loginForm.controls
  }
  ngOnInit(): void {
  }
  submit() {
    this.loginService.login(this.loginForm.value as LoginRequest).subscribe({
      next: (response: LoginResponse) => {
        if (response.statusCode === 200 && response.token.length > 0) {
          localStorage.setItem('token', response.token);
          this.router.navigateByUrl("/dashboard")
        } else {
          this.toasterService.show("Error", "Invalid credential")
        }
      }
    })
  }
}

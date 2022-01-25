import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MenuComponent } from './components/menu/menu.component';
import { HeaderComponent } from './components/header/header.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { ClientListComponent } from './components/client-list/client-list.component';
import { ClientAddComponent } from './components/client-add/client-add.component';
import { ClientUpdateComponent } from './components/client-update/client-update.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatTableModule } from '@angular/material/table';
import { MatSidenavModule } from '@angular/material/sidenav';
import { HomeComponent } from './components/home/home.component';
import { MatMenuModule } from '@angular/material/menu';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { ClientService } from './services/client-service';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { ReactiveFormsModule } from '@angular/forms';
import { HsnCodeListComponent } from './components/hsn-code-list/hsn-code-list.component';
import { HsnCodeAddComponent } from './components/hsn-code-add/hsn-code-add.component';
import { HsnCodeUpdateComponent } from './components/hsn-code-update/hsn-code-update.component';
import { HsnCodeService } from './services/hsn-code-service';
import { MatDialogModule } from '@angular/material/dialog';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { TokenInterceptor } from './interceptors/token-interceptor';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { SpinnerComponent } from './components/spinner/spinner.component';
import { ToastrModule } from 'ngx-toastr';
import { DeleteConfirmationComponent } from './components/delete-confirmation/delete-confirmation.component';
import { PaginationComponent } from './components/pagination/pagination.component';
import { VendorListComponent } from './components/vendor/vendor-list/vendor-list.component';
import { VendorAddComponent } from './components/vendor/vendor-add/vendor-add.component';
import { VendorUpdateComponent } from './components/vendor/vendor-update/vendor-update.component';
import { VendorService } from './services/vendor-service';
import { TransporterListComponent } from './components/transporter/transporter-list/transporter-list.component';
import { TransporterAddComponent } from './components/transporter/transporter-add/transporter-add.component';
import { TransporterUpdateComponent } from './components/transporter/transporter-update/transporter-update.component';
import { PurchaseListComponent } from './components/purchase/purchase-list/purchase-list.component';
import { PurchaseAddComponent } from './components/purchase/purchase-add/purchase-add.component';
import { PurchaseUpdateComponent } from './components/purchase/purchase-update/purchase-update.component';
import { PurchaseService } from './services/purchase-service';
import { SalesService } from './services/sales-service';
import { BailService } from './services/bail-service';
import { QualityService } from './services/quality-serice';
import { SalesListComponent } from './components/sales/sales-list/sales-list.component';
import { SalesAddComponent } from './components/sales/sales-add/sales-add.component';
import { SalesUpdateComponent } from './components/sales/sales-update/sales-update.component';

@NgModule({
  declarations: [
    AppComponent,
    MenuComponent,
    HeaderComponent,
    DashboardComponent,
    ClientListComponent,
    ClientAddComponent,
    ClientUpdateComponent,
    HomeComponent,
    HsnCodeListComponent,
    HsnCodeAddComponent,
    HsnCodeUpdateComponent,
    SpinnerComponent,
    DeleteConfirmationComponent,
    PaginationComponent,
    VendorListComponent,
    VendorAddComponent,
    VendorUpdateComponent,
    TransporterListComponent,
    TransporterAddComponent,
    TransporterUpdateComponent,
    PurchaseListComponent,
    PurchaseAddComponent,
    PurchaseUpdateComponent,
    SalesListComponent,
    SalesAddComponent,
    SalesUpdateComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    MatTableModule,
    MatSidenavModule,
    MatMenuModule,
    BrowserAnimationsModule,
    NgbModule,
    MatFormFieldModule,
    MatPaginatorModule,
    MatDividerModule,
    MatExpansionModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatAutocompleteModule,
    MatProgressSpinnerModule,
    ToastrModule.forRoot()
  ],
  providers: [ClientService, HsnCodeService,VendorService,PurchaseService, SalesService, BailService,QualityService,
    { provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true }],
  bootstrap: [AppComponent]
})
export class AppModule { }

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
import {MatTableModule} from '@angular/material/table';
import {MatSidenavModule} from '@angular/material/sidenav';
import { HomeComponent } from './components/home/home.component'; 
import {MatMenuModule} from '@angular/material/menu';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { ClientService } from './services/client-service';
import { HttpClientModule } from '@angular/common/http';
import {MatFormFieldModule} from '@angular/material/form-field'; 
import {MatPaginatorModule} from '@angular/material/paginator';
import { HsnCodeComponent } from './components/hsn-code/hsn-code.component'; 
import {MatDividerModule} from '@angular/material/divider';
import {MatExpansionModule} from '@angular/material/expansion';
import { ReactiveFormsModule } from '@angular/forms';
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
    HsnCodeComponent
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
    ReactiveFormsModule
  ],
  providers: [ClientService],
  bootstrap: [AppComponent]
})
export class AppModule { }

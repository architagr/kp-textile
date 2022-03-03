import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { PurchaseListComponent } from './components/purchase/purchase-list/purchase-list.component';
import { PurchaseAddComponent } from './components/purchase/purchase-add/purchase-add.component';
import { SalesAddComponent } from './components/sales/sales-add/sales-add.component';
import { SalesUpdateComponent } from './components/sales/sales-update/sales-update.component';
import { SalesListComponent } from './components/sales/sales-list/sales-list.component';
import { AuthGuard } from './services/auth-guard';
import { LoginComponent } from './components/login/login.component';
import { SpinnerComponent } from './components/spinner/spinner.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { HeaderComponent } from './components/header/header.component';
import { MenuComponent } from './components/menu/menu.component';
import { AppComponent } from './app.component';

export const declaration = [
  AppComponent,
  MenuComponent,
  HeaderComponent,
  DashboardComponent,
  HomeComponent,
  SpinnerComponent,
  PurchaseListComponent,
  PurchaseAddComponent,
  SalesListComponent,
  SalesAddComponent,
  SalesUpdateComponent,
  LoginComponent
]
const routes: Routes = [  
  {
    path: 'login',
    component: LoginComponent,
    pathMatch: 'full'
  },
  {
    path: 'master-data',
    loadChildren: () => import('./components/master-data-module/master-data.module').then(m => m.MasterDataModule),
    canActivate: [AuthGuard],
    pathMatch: 'prefix'
  },
  {
    path: "dashboard",
    component: HomeComponent,
  },
  
  {
    path: 'purchase',
    component: PurchaseListComponent,
  },
  {
    path: 'addpurchase',
    component: PurchaseAddComponent,
  },
  {
    path: 'sales',
    component: SalesListComponent,
    pathMatch: 'full'
  },
  {
    path: 'updatesales/:salesBillNumber',
    component: SalesUpdateComponent,
  },
  {
    path: 'addsales',
    component: SalesAddComponent,
  },

  {
    path: '',
    redirectTo:"/login",
    pathMatch: 'full'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }

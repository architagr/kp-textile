import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientAddComponent } from './components/client-add/client-add.component';
import { ClientListComponent } from './components/client-list/client-list.component';
import { ClientUpdateComponent } from './components/client-update/client-update.component';
import { HomeComponent } from './components/home/home.component';
import { HsnCodeListComponent } from './components/hsn-code-list/hsn-code-list.component';
import { TransporterAddComponent } from './components/transporter/transporter-add/transporter-add.component';
import { TransporterListComponent } from './components/transporter/transporter-list/transporter-list.component';
import { TransporterUpdateComponent } from './components/transporter/transporter-update/transporter-update.component';
import { VendorAddComponent } from './components/vendor/vendor-add/vendor-add.component';
import { VendorListComponent } from './components/vendor/vendor-list/vendor-list.component';
import { VendorUpdateComponent } from './components/vendor/vendor-update/vendor-update.component';
import { PurchaseListComponent } from './components/purchase/purchase-list/purchase-list.component';
import { PurchaseUpdateComponent } from './components/purchase/purchase-update/purchase-update.component';
import { PurchaseAddComponent } from './components/purchase/purchase-add/purchase-add.component';
import { SalesAddComponent } from './components/sales/sales-add/sales-add.component';
import { SalesUpdateComponent } from './components/sales/sales-update/sales-update.component';
import { SalesListComponent } from './components/sales/sales-list/sales-list.component';
import { QualityListComponent } from './components/quality/quality-list/quality-list.component';
import { AuthGuard } from './services/auth-guard';
import { LoginComponent } from './components/login/login.component';

const routes: Routes = [
  {
    path:'',
    component: LoginComponent,
    outlet:'loginRoute'
  },
  {
    path:'login',
    component: LoginComponent,
    outlet:'loginRoute'
  },
  {
    path:'master-data',
    loadChildren:()=>import ('./components/master-data-module/master-data.module').then(m=>m.MasterDataModule),
    canActivate: [AuthGuard],
    outlet:'nonLoginRoute'
  },
  {
    path: "dashboard",
    component: HomeComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "client",
    component: ClientListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "addclient",
    component: ClientAddComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "updateclient/:clientId",
    component: ClientUpdateComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "vendor",
    component: VendorListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "addvendor",
    component: VendorAddComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "updatevendor/:vendorId",
    component: VendorUpdateComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "transpoter",
    component: TransporterListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "addtranspoter",
    component: TransporterAddComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "updatetranspoter/:transpoterId",
    component: TransporterUpdateComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: "hsncode",
    component: HsnCodeListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'purchase',
    component: PurchaseListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'updatepurchase/:purchaseBillNumber',
    component: PurchaseUpdateComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'addpurchase',
    component: PurchaseAddComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'sales',
    component: SalesListComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'updatesales/:salesBillNumber',
    component: SalesUpdateComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'addsales',
    component: SalesAddComponent,
    outlet:'nonLoginRoute'
  },
  {
    path: 'quality',
    component: QualityListComponent,
    outlet:'nonLoginRoute'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {useHash: true})],
  exports: [RouterModule]
})
export class AppRoutingModule { }

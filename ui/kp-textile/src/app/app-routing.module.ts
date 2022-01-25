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

const routes: Routes = [
  {
    path: "",
    component: HomeComponent
  },
  {
    path: "dashboard",
    component: HomeComponent
  },
  {
    path: "client",
    component: ClientListComponent
  },
  {
    path: "addclient",
    component: ClientAddComponent
  },
  {
    path: "updateclient/:clientId",
    component: ClientUpdateComponent
  },
  {
    path: "vendor",
    component: VendorListComponent
  },
  {
    path: "addvendor",
    component: VendorAddComponent
  },
  {
    path: "updatevendor/:vendorId",
    component: VendorUpdateComponent
  },
  {
    path: "transpoter",
    component: TransporterListComponent
  },
  {
    path: "addtranspoter",
    component: TransporterAddComponent
  },
  {
    path: "updatetranspoter/:transpoterId",
    component: TransporterUpdateComponent
  },
  {
    path: "hsncode",
    component: HsnCodeListComponent
  },
  {
    path: 'purchase',
    component: PurchaseListComponent
  },
  {
    path: 'updatepurchase/:purchaseBillNumber',
    component: PurchaseUpdateComponent
  },
  {
    path: 'addpurchase',
    component: PurchaseAddComponent
  },
  {
    path: 'sales',
    component: SalesListComponent
  },
  {
    path: 'updatesales/:salesBillNumber',
    component: SalesUpdateComponent
  },
  {
    path: 'addsales',
    component: SalesAddComponent
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

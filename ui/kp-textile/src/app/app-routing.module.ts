import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientAddComponent } from './components/client-add/client-add.component';
import { ClientListComponent } from './components/client-list/client-list.component';
import { ClientUpdateComponent } from './components/client-update/client-update.component';
import { HomeComponent } from './components/home/home.component';
import { HsnCodeAddComponent } from './components/hsn-code-add/hsn-code-add.component';
import { HsnCodeListComponent } from './components/hsn-code-list/hsn-code-list.component';
import { SalesComponent } from './components/sales/sales.component';
import { VendorAddComponent } from './components/vendor/vendor-add/vendor-add.component';
import { VendorListComponent } from './components/vendor/vendor-list/vendor-list.component';
import { VendorUpdateComponent } from './components/vendor/vendor-update/vendor-update.component';

const routes: Routes = [
  {
    path:"",
    component:HomeComponent
  },
  {
    path:"dashboard",
    component:HomeComponent
  },
  {
    path:"client",
    component:ClientListComponent
  },
  {
    path:"addclient",
    component:ClientAddComponent
  },
  {
    path:"updateclient/:clientId",
    component:ClientUpdateComponent
  },
  {
    path:"vendor",
    component:VendorListComponent
  },
  {
    path:"addvendor",
    component:VendorAddComponent
  },
  {
    path:"updatevendor/:vendorId",
    component:VendorUpdateComponent
  },
  {
    path:"hsncode",
    component: HsnCodeListComponent
  },
  {
    path:"addhsncode",
    component: HsnCodeAddComponent
  },
  {
    path:'sales',
    component: SalesComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientAddComponent } from './components/client/client-add/client-add.component';
import { ClientListComponent } from './components/client/client-list/client-list.component';
import { ClientUpdateComponent } from './components/client/client-update/client-update.component';
import { GodownListComponent } from './components/godown/godown-list/godown-list.component';
import { HsnCodeListComponent } from './components/hsn-code/hsn-code-list/hsn-code-list.component';
import { QualityListComponent } from './components/quality/quality-list/quality-list.component';
import { TransporterAddComponent } from './components/transporter/transporter-add/transporter-add.component';
import { TransporterListComponent } from './components/transporter/transporter-list/transporter-list.component';
import { TransporterUpdateComponent } from './components/transporter/transporter-update/transporter-update.component';
import { VendorAddComponent } from './components/vendor/vendor-add/vendor-add.component';
import { VendorListComponent } from './components/vendor/vendor-list/vendor-list.component';
import { VendorUpdateComponent } from './components/vendor/vendor-update/vendor-update.component';

const routes: Routes = [

  {
    path: 'godown',
    component: GodownListComponent
  },
  {
    path: 'quality',
    component: QualityListComponent,
  },
  {
    path: "client",
    component: ClientListComponent,
  },
  {
    path: "addclient",
    component: ClientAddComponent,
  },
  {
    path: "updateclient/:clientId",
    component: ClientUpdateComponent,
  },
  {
    path: "vendor",
    component: VendorListComponent,
  },
  {
    path: "addvendor",
    component: VendorAddComponent,
  },
  {
    path: "updatevendor/:vendorId",
    component: VendorUpdateComponent,
  },
  {
    path: "transpoter",
    component: TransporterListComponent,
  },
  {
    path: "addtranspoter",
    component: TransporterAddComponent,
  },
  {
    path: "updatetranspoter/:transpoterId",
    component: TransporterUpdateComponent,
  },
  {
    path: "hsncode",
    component: HsnCodeListComponent,
  },
  {
    path: '',
    redirectTo:'/godown',
    pathMatch:"full"
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MasterDataRoutingModule { }

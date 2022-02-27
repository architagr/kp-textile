import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { GodownListComponent } from './components/godown/godown-list/godown-list.component';
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
  // {
  //   path: "hsncode",
  //   component: HsnCodeListComponent,
  // },
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

import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { SharedModule } from "src/app/shared-module/shared.module";
import { BailInfoComponent } from "./components/bail-info/bail-info.component";
import { ClientAddComponent } from "./components/client/client-add/client-add.component";
import { ClientListComponent } from "./components/client/client-list/client-list.component";
import { ClientUpdateComponent } from "./components/client/client-update/client-update.component";
import { AddGodownComponent } from "./components/godown/add-godown/add-godown.component";
import { GodownListComponent } from "./components/godown/godown-list/godown-list.component";
import { HsnCodeAddComponent } from "./components/hsn-code/hsn-code-add/hsn-code-add.component";
import { HsnCodeListComponent } from "./components/hsn-code/hsn-code-list/hsn-code-list.component";
import { AddProductComponent } from "./components/quality/add-product/add-product.component";
import { QualityAddComponent } from "./components/quality/quality-add/quality-add.component";
import { QualityListComponent } from "./components/quality/quality-list/quality-list.component";
import { QualityUpdateComponent } from "./components/quality/quality-update/quality-update.component";
import { TransporterAddComponent } from "./components/transporter/transporter-add/transporter-add.component";
import { TransporterListComponent } from "./components/transporter/transporter-list/transporter-list.component";
import { TransporterUpdateComponent } from "./components/transporter/transporter-update/transporter-update.component";
import { VendorAddComponent } from "./components/vendor/vendor-add/vendor-add.component";
import { VendorListComponent } from "./components/vendor/vendor-list/vendor-list.component";
import { VendorUpdateComponent } from "./components/vendor/vendor-update/vendor-update.component";
import { MasterDataRoutingModule } from "./master-data-routing.module";

@NgModule({
    declarations:[
        GodownListComponent,
        AddGodownComponent,
        QualityListComponent,
        QualityAddComponent,
        QualityUpdateComponent,
        AddProductComponent,
        BailInfoComponent,
        VendorAddComponent,
        VendorListComponent,
        VendorUpdateComponent,
        TransporterAddComponent,
        TransporterUpdateComponent,
        TransporterListComponent,
        HsnCodeAddComponent,
        HsnCodeListComponent,
        ClientAddComponent,
        ClientListComponent,
        ClientUpdateComponent
    ],
    imports: [
        CommonModule,
        SharedModule,
        MasterDataRoutingModule
    ],
    bootstrap: [GodownListComponent]
})
export class MasterDataModule { }

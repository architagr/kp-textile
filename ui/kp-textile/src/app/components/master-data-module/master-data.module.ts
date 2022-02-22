import { NgModule } from "@angular/core";
import { SharedModule } from "src/app/shared-module/shared.module";
import { BailInfoComponent } from "./components/bail-info/bail-info.component";
import { AddGodownComponent } from "./components/godown/add-godown/add-godown.component";
import { GodownListComponent } from "./components/godown/godown-list/godown-list.component";
import { AddProductComponent } from "./components/quality/add-product/add-product.component";
import { QualityAddComponent } from "./components/quality/quality-add/quality-add.component";
import { QualityListComponent } from "./components/quality/quality-list/quality-list.component";
import { QualityUpdateComponent } from "./components/quality/quality-update/quality-update.component";
import { MasterDataRoutingModule } from "./master-data-routing.module";
import { GodownService } from "./services/godown-service";
import { QualityService } from "./services/quality-serice";

@NgModule({
    declarations:[
        GodownListComponent,
        AddGodownComponent,
        QualityListComponent,
        QualityAddComponent,
        QualityUpdateComponent,
        AddProductComponent,
        BailInfoComponent
    ],
    imports: [
        SharedModule,
        MasterDataRoutingModule
    ],
    providers: [GodownService, QualityService],
    bootstrap: [GodownListComponent]
})
export class MasterDataModule { }

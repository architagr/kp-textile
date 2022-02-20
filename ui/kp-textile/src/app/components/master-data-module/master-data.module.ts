import { NgModule } from "@angular/core";
import { SharedModule } from "src/app/shared-module/shared.module";
import { AddGodownComponent } from "./components/godown/add-godown/add-godown.component";
import { GodownListComponent } from "./components/godown/godown-list/godown-list.component";
import { MasterDataRoutingModule } from "./master-data-routing.module";
import { GodownService } from "./services/godown-service";

@NgModule({
    declarations:[
        GodownListComponent,
        AddGodownComponent,
    ],
    imports: [
        SharedModule,
        MasterDataRoutingModule
    ],
    providers: [GodownService],
    bootstrap: [GodownListComponent]
})
export class MasterDataModule { }

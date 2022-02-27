import { ModuleWithProviders, NgModule } from "@angular/core";
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatTableModule } from '@angular/material/table';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatMenuModule } from '@angular/material/menu';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatDividerModule } from '@angular/material/divider';
import { MatExpansionModule } from '@angular/material/expansion';
import { ReactiveFormsModule } from '@angular/forms';
import { MatDialogModule } from '@angular/material/dialog';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from "@angular/common";
import { PaginationComponent } from "../components/pagination/pagination.component";
import { DeleteConfirmationComponent } from "../components/delete-confirmation/delete-confirmation.component";
import { BailService } from "../services/bail-service";
import { ClientService } from "../services/client-service";
import { DocumentService } from "../services/document-service";
// import { HsnCodeService } from "../services/hsn-code-service";
import { PurchaseService } from "../services/purchase-service";
import { SalesService } from "../services/sales-service";
import { ToastService } from "../services/toast-service";
import { TransporterService } from "../services/transporter-service";
import { VendorService } from "../services/vendor-service";
import { GodownService } from "../services/godown-service";
import { QualityService } from "../services/quality-serice";

const declaration = [
    PaginationComponent,
    DeleteConfirmationComponent
];

const imports = [
    MatTableModule,
    MatSidenavModule,
    MatMenuModule,
    // BrowserAnimationsModule,
    NgbModule,
    MatFormFieldModule,
    MatPaginatorModule,
    MatDividerModule,
    MatExpansionModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatAutocompleteModule,
    MatProgressSpinnerModule,
    MatIconModule,
    CommonModule
]
@NgModule({
    declarations: declaration,
    imports: imports,
    exports: [
        ...imports,
        ...declaration
    ]
})
export class SharedModule {
    static forRoot(): ModuleWithProviders<SharedModule> {
        return {
          ngModule: SharedModule,
          providers: [ BailService, 
            ClientService, 
            DocumentService, 
            // HsnCodeService, 
            PurchaseService, 
            SalesService, 
            ToastService, 
            TransporterService, 
            VendorService,
            GodownService,
            QualityService
         ],
        }
      }
}
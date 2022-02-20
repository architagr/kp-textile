import { NgModule } from "@angular/core";
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


@NgModule({
    imports: [
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
    ],
    exports: [
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
})
export class SharedModule{

}
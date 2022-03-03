import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ToastrService } from 'ngx-toastr';
import { VendorDto } from 'src/app/models/vendor-models';
import { VendorService } from 'src/app/services/vendor-service';
import { DeleteConfirmationComponent } from '../../../../delete-confirmation/delete-confirmation.component';


@Component({
  selector: 'app-vendor-list',
  templateUrl: './vendor-list.component.html',
  styleUrls: ['./vendor-list.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*', })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
    trigger('detailExpand1', [
      state('collapsed', style({ 'padding-top': '0px', 'padding-bottom': '0px' })),
      state('expanded', style({ 'padding-top': '*', 'padding-bottom': '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class VendorListComponent implements OnInit {
  filterForm: FormGroup
  searchText: string = ''
  displayedColumns: string[] = ['CompanyName', 'PaymentTerms', 'Status', 'ContactInfo', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 10;
  total: number = 0;
  lastEvalutionKey: any = null
  vendors: VendorDto[] = []
  vendorsAll: VendorDto[] = []
  expandedElement: VendorDto | null = null
  constructor(
    public vendorService: VendorService,
    private toastr: ToastrService,
    public dialog: MatDialog
  ) { 
    this.filterForm = new FormGroup({
      searchText: new FormControl('')
    });

  }

  ngOnInit(): void {
    this.getVendors();
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getVendors();
  }

  get startIndex(): number{
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number{
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.vendorsAll.length) {
      endIndex = this.vendorsAll.length - 1
    }
    return endIndex
  }
  onNextPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;

    if (this.vendorsAll.length >= this.endIndex + 1) {
      this.getVendorFromLocalList()
    } else {
      this.getVendors();
    }
  }
  getVendorFromLocalList(){
    this.vendors = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex 
    for (; startIndex <= endIndex; startIndex++) {
      this.vendors.push(this.vendorsAll[startIndex])
    }
  }
  onPrevPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getVendorFromLocalList();
  }
  getVendors() {
    this.vendorService.getVendorList(this.pageSize, this.searchText, this.lastEvalutionKey).subscribe(data => {
      this.vendors = data.data
      this.addToAllVendorList(data.data);
      this.lastEvalutionKey = data.lastEvalutionKey
      this.total = data.total
    });
  }
  addToAllVendorList(data: VendorDto[]) {
    if(data)
    data.forEach(x => {
      if (!this.vendorsAll.some(y => y.vendorId === x.vendorId)) {
        this.vendorsAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getVendors();
  }
  deleteVendor(vendor: VendorDto) {
    const dialogRef = this.dialog.open(DeleteConfirmationComponent, {
      data: { heading: 'Delete Vendor', message: `Are you sure you want to delete ${vendor.companyName} ?` },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        this.delete(vendor.vendorId);
      }
    });
  }
  delete(vendorId: string) {
    this.vendorService.deleteVendor(vendorId).subscribe({
      next: (data) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Vendor deleted.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        this.getVendors();
      },
      error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Error in deleting vendor.', 'Error', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-danger alert-with-icon",
          positionClass: 'toast-top-right'
        });
      }
    })
  }
}

import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { MatDialog } from '@angular/material/dialog';
import { DeleteConfirmationComponent } from '../../delete-confirmation/delete-confirmation.component';
import { FormControl, FormGroup } from '@angular/forms';
import { InventoryDto } from 'src/app/models/item-model';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { SalesService } from 'src/app/services/sales-service';
import { QualityService } from 'src/app/services/quality-serice';
import { QualityDto } from 'src/app/models/quality-model';

@Component({
  selector: 'app-sales-list',
  templateUrl: './sales-list.component.html',
  styleUrls: ['./sales-list.component.scss'],
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
export class SalesListComponent implements OnInit {
  filterForm: FormGroup
  searchText: string = ''
  displayedColumns: string[] = ['SalesBillNumber', 'HsnCode', 'TotalQuantity', 'SalesDate', 'ChallanNo', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 10;
  total: number = 0;
  lastEvalutionKey: any = null
  sales: InventoryDto[] = []
  salesAll: InventoryDto[] = []
  expandedElement: InventoryDto | null = null
  hsnCodes: HnsCodeDto[] = [];
  qualities: QualityDto[] = [];

  constructor(
    public salesService: SalesService,
    private toastr: ToastrService,
    public dialog: MatDialog,
    private qualityService: QualityService,
    private hsnCodeService: HsnCodeService
  ) {
    this.filterForm = new FormGroup({
      searchText: new FormControl('')
    });

  }
  getTotalQuanity(purchaseOrder: InventoryDto): number {
    let total = 0;
    purchaseOrder.bailDetails.forEach(x => total = total + x.billedQuantity)
    return total
  }
  ngOnInit(): void {
    this.getAllCodes();
  }

  getAllCodes() {
    this.hsnCodeService.getAllHsnCode().subscribe({
      next: (data) => {
        this.hsnCodes = data.data
      },
      complete: () => {
        this.getAllQuality();
      }
    });
  }
  getAllQuality() {
    this.qualityService.getAllQualities().subscribe({
      next: (data) => {
        this.qualities = data.data
      },
      complete: () => {
        this.getSales()
      }
    });
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getSales();
  }

  get startIndex(): number {
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number {
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.salesAll.length) {
      endIndex = this.salesAll.length - 1
    }
    return endIndex
  }
  onNextPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;

    if (this.salesAll.length >= this.endIndex + 1) {
      this.getPurchaseFromLocalList()
    } else {
      this.getSales();
    }
  }
  getPurchaseFromLocalList() {
    this.sales = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex
    for (; startIndex <= endIndex; startIndex++) {
      this.sales.push(this.salesAll[startIndex])
    }
  }
  onPrevPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getPurchaseFromLocalList();
  }
  getSales() {

    this.salesService.getAllSales(this.lastEvalutionKey, this.pageSize).subscribe({
      next: (data) => {
        this.sales = data.data
        data.data.forEach(x => {
          x.hsnCode = this.hsnCodes.find(y => y.id === x.hsnCode)!.hnsCode;
          x.bailDetails.forEach(x => {
            x.quality = this.qualities.find(y => y.id === x.quality)!.name
          })
        })
        this.addToAllSalesList(data.data);
        this.lastEvalutionKey = data.lastEvalutionKey
        this.total = data.total;
      },
      error: (err) => {
        this.sales = [];
        this.lastEvalutionKey = null

      }
    });
  }
  addToAllSalesList(data: InventoryDto[]) {
    data.forEach(x => {
      if (!this.salesAll.some(y => y.billNo === x.billNo)) {
        this.salesAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getSales();
  }
  deleteSales(salesOrder: InventoryDto) {
    const dialogRef = this.dialog.open(DeleteConfirmationComponent, {
      data: { heading: 'Delete Sales Order', message: `Are you sure you want to delete ${salesOrder.billNo} ?` },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        this.delete(salesOrder.billNo);
      }
    });
  }
  delete(purchaseBillNo: string) {
    this.salesService.deleteSales(purchaseBillNo).subscribe({
      next: (data) => {
        this.getSales();
      },
      error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Error in deleting Sales order.', 'Error', {
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

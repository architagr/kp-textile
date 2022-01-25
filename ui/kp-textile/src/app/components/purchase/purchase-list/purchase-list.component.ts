import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { MatDialog } from '@angular/material/dialog';
import { DeleteConfirmationComponent } from '../../delete-confirmation/delete-confirmation.component';
import { FormControl, FormGroup } from '@angular/forms';
import { PurchaseService } from 'src/app/services/purchase-service';
import { InventoryDto } from 'src/app/models/item-model';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';

@Component({
  selector: 'app-purchase-list',
  templateUrl: './purchase-list.component.html',
  styleUrls: ['./purchase-list.component.scss'],
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
export class PurchaseListComponent implements OnInit {
  filterForm: FormGroup
  searchText: string = ''
  displayedColumns: string[] = ['PurchaseBillNumber', 'HsnCode', 'TotalQuantity', 'PurchaseDate', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 10;
  total: number = 0;
  lastEvalutionKey: any = null
  purchases: InventoryDto[] = []
  purchasesAll: InventoryDto[] = []
  expandedElement: InventoryDto | null = null
  hsnCodes: HnsCodeDto[] = [];
  constructor(
    public purchaseService: PurchaseService,
    private toastr: ToastrService,
    public dialog: MatDialog,
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
        this.getPurchases();
      }
    });
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getPurchases();
  }

  get startIndex(): number {
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number {
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.purchasesAll.length) {
      endIndex = this.purchasesAll.length - 1
    }
    return endIndex
  }
  onNextPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;

    if (this.purchasesAll.length >= this.endIndex + 1) {
      this.getPurchaseFromLocalList()
    } else {
      this.getPurchases();
    }
  }
  getPurchaseFromLocalList() {
    this.purchases = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex
    for (; startIndex <= endIndex; startIndex++) {
      this.purchases.push(this.purchasesAll[startIndex])
    }
  }
  onPrevPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getPurchaseFromLocalList();
  }
  getPurchases() {

    this.purchaseService.getAllPurchase(this.lastEvalutionKey, this.pageSize).subscribe({
      next: (data) => {
        this.purchases = data.data
        data.data.forEach(x => {
          x.hsnCode = this.hsnCodes.find(y => y.id === x.hsnCode)!.hnsCode;
        })
        this.addToAllPurchaseList(data.data);
        this.lastEvalutionKey = data.lastEvalutionKey
        this.total = data.total;
      },
     
    });
  }
  addToAllPurchaseList(data: InventoryDto[]) {
    data.forEach(x => {
      if (!this.purchasesAll.some(y => y.billNo === x.billNo)) {
        this.purchasesAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getPurchases();
  }
  deletePurchase(purchaseOrder: InventoryDto) {
    const dialogRef = this.dialog.open(DeleteConfirmationComponent, {
      data: { heading: 'Delete Purchase Order', message: `Are you sure you want to delete ${purchaseOrder.billNo} ?` },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        this.delete(purchaseOrder.billNo);
      }
    });
  }
  delete(purchaseBillNo: string) {
    this.purchaseService.deletePurchaseOrder(purchaseBillNo).subscribe({
      next: (data) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Purchase order deleted.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        this.getPurchases();
      },
      error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Error in deleting purchase order.', 'Error', {
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

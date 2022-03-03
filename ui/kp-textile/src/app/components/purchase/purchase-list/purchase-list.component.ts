import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { MatDialog } from '@angular/material/dialog';
import { DeleteConfirmationComponent } from '../../delete-confirmation/delete-confirmation.component';
import { FormControl, FormGroup } from '@angular/forms';
import { PurchaseService } from 'src/app/services/purchase-service';
import { InventoryDto, PurchaseDto } from 'src/app/models/item-model';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { QualityService } from 'src/app/services/quality-serice';
import { GodownService } from 'src/app/services/godown-service';
import { ProductDto, QualityDto } from 'src/app/models/quality-model';
import { VendorDto } from 'src/app/models/vendor-models';
import { GodownDto } from 'src/app/models/godown-model';
import { VendorService } from 'src/app/services/vendor-service';

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
  displayedColumns: string[] = ['purchaseBillNo', 'vendorName', 'productName', 'qualityName', 'hsnCode', 'purchaseStatus', 'date'];
  pageNumber: number = 0;
  pageSize: number = 1;
  total: number = 0;
  lastEvalutionKey: any = null
  purchases: PurchaseDto[] = []
  purchasesAll: PurchaseDto[] = []
  expandedElement: PurchaseDto | null = null
  qualities: QualityDto[] = [];
  products: ProductDto[] = [];
  vendors: VendorDto[] = [];
  godowns: GodownDto[] = [];
  selectedGodown: string = '';
  constructor(
    public purchaseService: PurchaseService,
    private toastr: ToastrService,
    public dialog: MatDialog,
    private qualityService: QualityService,
    private godownService: GodownService,
    private vendorService: VendorService
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
    this.getAllGodowns();
    this.getAllQuality();
    this.getAllProducts();
    this.getAllVendors();
  }

  getAllGodowns() {
    this.godownService.getAllGodown().subscribe({
      next: (data) => {
        this.godowns = data.data
      },
    });
  }
  getAllQuality() {
    this.qualityService.getAllQualities().subscribe({
      next: (data) => {
        this.qualities = data.data
      },
    });
  }
  getAllProducts() {
    this.qualityService.getAllProduct().subscribe({
      next: (data) => {
        this.products = data.data
      },
    });
  }
  getAllVendors() {
    this.vendorService.getAllVendors().subscribe({
      next: (data: VendorDto[]) => {
        this.vendors = data;
      }
    })
  }
  selectedIndexChange(event: number){
    this.purchases = [];
    this.purchasesAll = [];
    this.lastEvalutionKey = null;
    this.selectedGodown = this.godowns[event].id;
    this.getPurchases();
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
    if(this.selectedGodown)
    this.purchaseService.getAllPurchase(this.lastEvalutionKey, this.pageSize, this.selectedGodown).subscribe({
      next: (data) => {
        
        data.data.forEach(x => {
          x.vendorName = this.vendors.find(y => y.vendorId === x.vendorId)!.companyName;
          x.productName= this.products.find(y => y.id === x.productId)!.name;
          let quality = this.qualities.find(y => y.id === x.qualityId);
          x.qualityName = quality!.name;
          x.hsnCode = quality!.hsnCode;
        })
        this.purchases = data.data;
        this.addToAllPurchaseList(data.data);
        this.lastEvalutionKey = data.lastEvalutionKey
        this.total = data.total;
      },

    });
  }
  addToAllPurchaseList(data: PurchaseDto[]) {
    data.forEach(x => {
      if (!this.purchasesAll.some(y => y.purchaseBillNo === x.purchaseBillNo)) {
        this.purchasesAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getPurchases();
  }
}

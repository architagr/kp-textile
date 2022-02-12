import { Component, OnInit } from '@angular/core';
import { BailInfoResponse } from 'src/app/models/item-model';
import { ProductDto, ProductListResponse, QualityDto, QualityListItemDto, QualityListResponse } from 'src/app/models/quality-model';
import { BailService } from 'src/app/services/bail-service';
import { QualityService } from 'src/app/services/quality-serice';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { catchError, from, mergeMap, of, toArray } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { QualityAddComponent } from '../quality-add/quality-add.component';
import { BailInfoComponent } from '../../bail-info/bail-info.component';
import { AddProductComponent } from '../add-product/add-product.component';
import { ToastService } from 'src/app/services/toast-service';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-quality-list',
  templateUrl: './quality-list.component.html',
  styleUrls: ['./quality-list.component.scss'],
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
export class QualityListComponent implements OnInit {
  filterForm: FormGroup
  allQuality: QualityListItemDto[] = [];
  allProduct: ProductDto[] = [];
  qualities: QualityListItemDto[] = [];
  expandedElement: QualityListItemDto | null = null
  displayedColumns: string[] = ['QualityName', 'ProductName', 'HsnCode', 'RemainingQuantity', 'NoBale'];
  pageSize = 5;
  pageNumber = 0;
  searchText = ''
  constructor(
    private qualityService: QualityService,
    private bailService: BailService,
    public dialog: MatDialog,
    private toastService: ToastService
  ) {

    this.filterForm = new FormGroup({
      searchText: new FormControl('')
    })
  }
  
  ngOnInit(): void {
    this.getAllProduct(true);
  }
  getAllProduct(updateQualities: boolean) {
    this.allProduct = [];
    this.qualityService.getAllProduct().subscribe({
      next: (response: ProductListResponse) => {
        this.allProduct = response.data;
        if (updateQualities)
          this.getAllQuality();
      }
    })
  }
  getAllQuality() {
    this.allQuality = [];
    this.pageNumber = 0;
    this.qualityService.getAllQualities().subscribe({
      next: (response: QualityListResponse) => {
        response.data.forEach(x => {
          x.productName = this.allProduct.find(y => y.id === x.productId)!.name
        })
        this.getAllSalableBails(response.data);
      },
      error: (err) => {
        this.toastService.show("Error", err.error?.errorMessage ?? err.message)
      }
    });
  }
  getAllSalableBails(qualities: QualityDto[]) {
    from(qualities).pipe(mergeMap((element) =>
      this.bailService.getBailInfoByQuality(element.id).pipe(catchError(error => {
        return of({} as BailInfoResponse)
      }),
      ),
      this.pageSize),
      toArray()).subscribe({
        next: (response) => {
          for (let index = 0; index < qualities.length; index++) {
            const element = qualities[index];
            const data = response.find(x => x.statusCode === 200 && x.purchase.length > 0 && x.purchase[0].quality === element.id);

            if (data && data!.purchase!.length > 0) {
              let total = 0;
              data!.purchase.forEach(x => total = total + x.pendingQuantity);

              this.allQuality.push({ ...element, pendingQuantity: total, bailDetails: data!.purchase } as QualityListItemDto);
            } else {
              this.allQuality.push({ ...element, pendingQuantity: 0, bailDetails: [] } as QualityListItemDto);
            }
          }
          this.getQualitiesFromLocalList();
        }
      })

  }

  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.pageNumber = 0;
  }

  get startIndex(): number {
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number {
    let endIndex = this.startIndex + this.pageSize - 1
    let qualities = this.filterAllQuality()
    if (endIndex >= qualities.length) {
      endIndex = qualities.length - 1
    }
    return endIndex
  }
  filterAllQuality(): QualityListItemDto[]{
    let qualities = this.allQuality
    if(this.searchText.length>0){
      qualities = this.allQuality.filter(x=>x.name.toLocaleLowerCase().indexOf(this.searchText)!==-1
      || x.hsnCode.toLocaleLowerCase().indexOf(this.searchText)!==-1
      || x.hsnCode.toLocaleLowerCase().indexOf(this.searchText)!==-1) 
    }
    return qualities;
  }
  onPageChange(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getQualitiesFromLocalList();
  }
  
  applyFilter(){
    this.pageNumber = 0;
    this.searchText = this.filterForm.value['searchText']
    this.getQualitiesFromLocalList();
  }
  getQualitiesFromLocalList() {
    this.qualities = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex
    console.log(`endIndex: ${endIndex}`)
    let qualities = this.filterAllQuality()

    for (; startIndex <= endIndex; startIndex++) {
      this.qualities.push(qualities[startIndex])
    }
  }
  addProductOpenDialog() {
    const dialogRef = this.dialog.open(AddProductComponent, {
      data: { id: '', name: '' } as ProductDto,
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        const product = (result as ProductDto)
        if (product.id === "") {
          this.addProduct(product);
        }
      }
    });
  }

  addProduct(product: ProductDto) {
    this.qualityService.addProduct(product.name).subscribe(response => {
      console.log(`data added `, response)
      this.getAllProduct(false);
    })
  }

  addQualityOpenDialog(): void {
    const dialogRef = this.dialog.open(QualityAddComponent, {
      data: {
        quality: {
          id: '',
          name: '',
          hsnCode: '',
          productId: '',
          productName: ''
        } as QualityDto,
        products: this.allProduct
      },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        const hsnCode = (result as QualityDto)
        if (hsnCode.id === "") {
          this.addQuality(hsnCode);
        }
      }
    });
  }

  addQuality(quality: QualityDto) {
    this.qualityService.addQuality(quality).subscribe(response => {
      console.log(`data added `, response)
      this.getAllQuality();
    })
  }

  showBaleInfo(baleNumber: string) {

    this.bailService.getBailInfo(baleNumber).subscribe({
      next: (response) => {
        this.openBaleInfo(baleNumber, response);
      },
      error: (err) => {
        console.log(`error in getting bale info `, err)
      }
    })
  }

  openBaleInfo(baleNumber: string, info: BailInfoResponse) {
    const dialogRef = this.dialog.open(BailInfoComponent, {
      data: { info: info, baleName: baleNumber },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);

    });
  }
}

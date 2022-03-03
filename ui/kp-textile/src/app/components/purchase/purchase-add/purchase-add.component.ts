import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { empty, EMPTY, expand, map, reduce } from 'rxjs';
import { GodownDto } from 'src/app/models/godown-model';
import { AddPurchaseDataRequest, InventoryDto } from 'src/app/models/item-model';
import { ProductDto, QualityDto } from 'src/app/models/quality-model';
import { VendorDto, VendorListResponse } from 'src/app/models/vendor-models';
import { GodownService } from 'src/app/services/godown-service';
import { PurchaseService } from 'src/app/services/purchase-service';
import { QualityService } from 'src/app/services/quality-serice';
import { ToastService } from 'src/app/services/toast-service';
import { VendorService } from 'src/app/services/vendor-service';
@Component({
  selector: 'app-purchase-add',
  templateUrl: './purchase-add.component.html',
  styleUrls: ['./purchase-add.component.scss']
})
export class PurchaseAddComponent implements OnInit {
  addPurchaseForm: FormGroup;
  hsnCode: string = "";
  qualities: QualityDto[] = [];
  uiQualities: QualityDto[] = [];
  products: ProductDto[] = [];
  vendors: VendorDto[] = [];
  godowns: GodownDto[] = [];
  constructor(
    private router: Router,
    private toastr: ToastService,
    private purchaseService: PurchaseService,
    private fb: FormBuilder,
    private qualityService: QualityService,
    private godownService: GodownService,
    private vendorService: VendorService
  ) {
    this.addPurchaseForm = this.fb.group({
      purchaseDetails: this.fb.group({
        godownId: new FormControl('', [Validators.required]),
        purchaseBillNo: new FormControl('', [Validators.required]),
        vendorId: new FormControl('', [Validators.required]),
        date: new FormControl(new Date(), [Validators.required]),
        productId: new FormControl('', [Validators.required]),
        qualityId: new FormControl('', [Validators.required]),
      }),
      baleDetails: this.fb.array([
        this.getBailFormGroup()
      ], [Validators.required]),
    })
    this.purchaseDetailsFC['productId'].valueChanges.subscribe(x => {
      this.uiQualities = this.qualities.filter(qual => qual.productId === x)
    })
    this.purchaseDetailsFC['qualityId'].valueChanges.subscribe(x => {
      this.hsnCode = this.qualities.find(qual => qual.id === x)?.hsnCode!
    })
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
  

  private getBailFormGroup(): FormGroup {
    return this.fb.group({
      baleNo: new FormControl('', [Validators.required]),
      billedQuantity: new FormControl(0, [Validators.required, Validators.min(1)]),
    });
  }
  get formControls(): { [key: string]: AbstractControl } {
    return this.addPurchaseForm.controls
  }
  get purchaseDetailsFC(): { [key: string]: AbstractControl } {
    return (this.formControls['purchaseDetails'] as FormGroup).controls
  }
  get baleDetails(): FormArray {
    return this.addPurchaseForm.controls['baleDetails'] as FormArray
  }
  addBale() {
    this.baleDetails.push(this.getBailFormGroup());
  }
  removeBale(removeIndex: number) {
    this.baleDetails.removeAt(removeIndex);
  }
  getBaleControl(index: number): { [key: string]: AbstractControl } {
    return (this.baleDetails.controls[index] as FormGroup).controls
  }
  ngOnInit(): void {
    this.getAllGodowns();
    this.getAllProducts();
    this.getAllQuality();
    this.getAllVendors();
  }
  submitData() {
    let data: AddPurchaseDataRequest = this.addPurchaseForm.value
    data.purchaseDetails.date = new Date(data.purchaseDetails.date)
    this.purchaseService.addPurchase(data).subscribe({
      next: (data) => {
        this.toastr.show("Success", "Purchase order added.")
        setTimeout(() => {
          this.router.navigate(['/purchase']);
        }, 2000);
      },
      error: (err) => {
        debugger;
        let errorMessage = `${err.error.errors[0].errorMessage}`
        this.toastr.show("Error", `${err.error.errorMessage}.<br/>${errorMessage}`)
      },
      complete: () => {
        console.log(`complete`)
      }
    })
  }
}

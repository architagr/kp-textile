import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { InventoryDto } from 'src/app/models/item-model';
import { QualityDto } from 'src/app/models/quality-model';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { PurchaseService } from 'src/app/services/purchase-service';
import { QualityService } from 'src/app/components/master-data-module/services/quality-serice';
@Component({
  selector: 'app-purchase-add',
  templateUrl: './purchase-add.component.html',
  styleUrls: ['./purchase-add.component.scss']
})
export class PurchaseAddComponent implements OnInit {
  addPurchaseForm: FormGroup;
  showSpinnerCount= 0;
  hsnCodes: HnsCodeDto[] = [];
  qualities: QualityDto[] = [];
  constructor(
    private router: Router,
    private toastr: ToastrService,
    private purchaseService: PurchaseService,
    private fb: FormBuilder,
    private hsnCodeService: HsnCodeService,
    private qualityService: QualityService,
  ) {
    this.addPurchaseForm = this.fb.group({
      billNo: new FormControl('', [Validators.required]),
      hsnCode: new FormControl('', [Validators.required]),
      purchaseDate: new FormControl(new Date(), [Validators.required]),
      bailDetails: this.fb.array([
        this.getBailFormGroup()
      ], [Validators.required]),
    })

  }
  getAllQuality() {
    this.showSpinnerCount++;
    this.qualityService.getAllQualities().subscribe({
      next:(data)=>{
          this.qualities = data.data
          this.showSpinnerCount--;
      },
      complete:()=>{

      }
    });
  }

  getAllCodes() {
    this.showSpinnerCount++;
    this.hsnCodeService.getAllHsnCode().subscribe({
      next:(data)=>{
          this.hsnCodes = data.data
          this.showSpinnerCount--;
      },
      complete:()=>{

      }
    });
  }
  private getBailFormGroup(): FormGroup {
    return this.fb.group({
      bailNo: new FormControl('', [Validators.required]),
      quality: new FormControl('', [Validators.required]),
      billedQuantity: new FormControl(0, [Validators.required, Validators.min(1)]),
      receivedQuantity: new FormControl(0),
    });
  }
  get formControls(): { [key: string]: AbstractControl } {
    return this.addPurchaseForm.controls
  }
  get bailDetails(): FormArray {
    return this.addPurchaseForm.controls['bailDetails'] as FormArray
  }
  addBail() {
    this.bailDetails.push(this.getBailFormGroup());
  }
  removeBail(removeIndex: number) {
    this.bailDetails.removeAt(removeIndex);
  }
  getBailControl(index: number): { [key: string]: AbstractControl } {
    return (this.bailDetails.controls[index] as FormGroup).controls
  }
  ngOnInit(): void {
    this.getAllCodes();
    this.getAllQuality();
  }
  submitData(){
    let data: InventoryDto = this.addPurchaseForm.value
    data.purchaseDate = new Date(data.purchaseDate)
    this.purchaseService.addPurchase(data).subscribe({
      next:(data)=>{
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Purchase order added.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        setTimeout(() => {
          this.router.navigate(['/purchase']);
        }, 2000);
      },
      error:(err)=>{
        let errorMessage  = `${err.error.errors[0].errorMessage}`
        this.toastr.info(`<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> ${err.error.errorMessage}.<br/>${errorMessage}`, 'Error', {
          disableTimeOut: false,
          timeOut: 3000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-danger alert-with-icon",
          positionClass: 'toast-top-right'
        });
      },
      complete: () =>{
        console.log(`complete`)
      }
    })
  }
}

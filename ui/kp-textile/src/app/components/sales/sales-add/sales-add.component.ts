import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { BailDetailsDto, InventoryDto } from 'src/app/models/item-model';
import { QualityDto } from 'src/app/models/quality-model';
import { TransporterDto } from 'src/app/models/transporter-model';
import { BailService } from 'src/app/services/bail-service';
import { DocumentService } from 'src/app/services/document-service';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { QualityService } from 'src/app/services/quality-serice';
import { SalesService } from 'src/app/services/sales-service';
import { TransporterService } from 'src/app/services/transporter-service';

@Component({
  selector: 'app-sales-add',
  templateUrl: './sales-add.component.html',
  styleUrls: ['./sales-add.component.scss']
})
export class SalesAddComponent implements OnInit {
  addPurchaseForm: FormGroup;
  showSpinnerCount = 0;
  hsnCodes: HnsCodeDto[] = [];
  allTransporters: TransporterDto[] = [];
  qualities: QualityDto[] = [];
  overAllBailInfo: { [propName: string]: BailDetailsDto[] } = {};
  constructor(
    private router: Router,
    private toastr: ToastrService,
    private salesService: SalesService,
    private fb: FormBuilder,
    private hsnCodeService: HsnCodeService,
    private qualityService: QualityService,
    private transporterService: TransporterService,
    private bailService: BailService,
    private documentService: DocumentService
  ) {
    this.addPurchaseForm = this.fb.group({
      billNo: new FormControl('', [Validators.required]),
      hsnCode: new FormControl('', [Validators.required]),
      salesDate: new FormControl(new Date(), [Validators.required]),
      lrNo: new FormControl('', [Validators.required]),
      challanNo: new FormControl('', [Validators.required]),
      transporterId: new FormControl('', [Validators.required]),
      bailDetails: this.fb.array([
        this.getBailFormGroup()
      ], [Validators.required]),
    })

  }
  getAllQuality() {
    this.showSpinnerCount++;
    this.qualityService.getAllQualities().subscribe({
      next: (data) => {
        this.qualities = data.data
        this.showSpinnerCount--;
        this.getAllSalableBails();
      },
      complete: () => {

      }
    });
  }
  getAllSalableBails() {
    for (let index = 0; index < this.qualities.length; index++) {
      const element = this.qualities[index];
      this.showSpinnerCount++;
      this.bailService.getBailInfoByQuality(element.id).subscribe({
        next: (response) => {
          this.overAllBailInfo[element.id] = [];
          this.overAllBailInfo[element.id] = response.purchase;
          this.showSpinnerCount--;
        },
        error: () => {
          this.overAllBailInfo[element.id] = [];
          this.showSpinnerCount--;
        }
      })
    }
  }
  getAllTransporter() {
    this.showSpinnerCount++;
    this.transporterService.getAllTransporter(10, '', null).subscribe({
      next: (data) => {
        this.allTransporters = data.data;
        this.showSpinnerCount--;
      },

    })
  }
  getAllCodes() {
    this.showSpinnerCount++;
    this.hsnCodeService.getAllHsnCode().subscribe({
      next: (data) => {
        this.hsnCodes = data.data
        this.showSpinnerCount--;
      },
      complete: () => {

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
  getBailsDetails(index: number): BailDetailsDto[] {
    let qualityId = this.getBailControl(index)['quality'].value
    return this.overAllBailInfo[qualityId]
  }
  ngOnInit(): void {
    this.getAllCodes();
    this.getAllQuality();
    this.getAllTransporter();
  }
  submitData() {
    let finalData: InventoryDto = this.addPurchaseForm.value
    finalData.salesDate = new Date(finalData.salesDate)
    this.salesService.addSales(finalData).subscribe({
      next: (data) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Sales order added.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });

        this.documentService.getChallan(finalData).subscribe({
          next: (response) => {
            var w = window.open('about:blank');
            w!.document.open();
            w!.document.write(response);
            this.router.navigate(['/sales']);
          }
        })
        // setTimeout(() => {
        //   this.router.navigate(['/sales']);
        // }, 2000);
      },
      error: (err) => {
        let errorMessage = `${err.error.errors[0].errorMessage}`
        this.toastr.info(`<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> ${err.error.errorMessage}.<br/>${errorMessage}`, 'Error', {
          disableTimeOut: false,
          timeOut: 3000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-danger alert-with-icon",
          positionClass: 'toast-top-right'
        });
      },
      complete: () => {
        console.log(`complete`)
      }
    })
  }
}

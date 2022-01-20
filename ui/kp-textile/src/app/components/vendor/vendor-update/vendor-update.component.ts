import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { AddressType, PaymentTerm, PersonType, Status } from 'src/app/models/client-model';
import { ToastrService } from 'ngx-toastr';
import { VendorService } from 'src/app/services/vendor-service';

@Component({
  selector: 'app-vendor-update',
  templateUrl: './vendor-update.component.html',
  styleUrls: ['./vendor-update.component.scss']
})
export class VendorUpdateComponent implements OnInit {
  updateVendorForm: FormGroup;

  paymentTermsValues: string[] = [];
  statusValues: string[] = [];
  addressTypeValues: string[] = [];
  personTypeValues: string[] = [];

  private getContactInfoFormGroup(): FormGroup {
    return this.fb.group({
      email: new FormControl('', [Validators.email]),
      landline: new FormControl(''),
      mobile: new FormControl('', [Validators.pattern('^([0|\+[0-9]{1,5})?([7-9][0-9]{9})$')]),
      whatsapp: new FormControl('', [Validators.pattern('^([0|\+[0-9]{1,5})?([7-9][0-9]{9})$')]),
    })
  }
  private getAddressFormGroup(): FormGroup {
    return this.fb.group({
      label: new FormControl('', [Validators.required]),
      addressType: new FormControl('', [Validators.required]),
      addressLine1: new FormControl('', [Validators.required]),
      addressLine2: new FormControl(''),
      pincode: new FormControl('', [Validators.required]),
      country: new FormControl('', [Validators.required]),
      state: new FormControl('', [Validators.required]),
      city: new FormControl('', [Validators.required]),
      landline: new FormControl(''),
      mobile: new FormControl('', [Validators.pattern('^([0|\+[0-9]{1,5})?([7-9][0-9]{9})$')]),
    });
  }

  private getContactPersonFormGroup(): FormGroup {
    return this.fb.group({
      branchId: new FormControl(''),
      sortKey: new FormControl(''),
      vendorId: new FormControl(''),
      contactId: new FormControl(''),
      salutation: new FormControl('', [Validators.required]),
      firstName: new FormControl('', [Validators.required]),
      lastName: new FormControl(''),
      department: new FormControl(''),
      personType: new FormControl('', [Validators.required]),
      remark: new FormControl(''),
      address: this.getAddressFormGroup(),
      contactInfo: this.getContactInfoFormGroup(),
    })
  }
  constructor(
    private route: ActivatedRoute,
    private vendorService: VendorService,
    private fb: FormBuilder,
    private toastr: ToastrService
  ) {
    this.updateVendorForm = this.fb.group({
      branchId: new FormControl(''),
      sortKey: new FormControl(''),
      vendorId: new FormControl(''),
      companyName: new FormControl('', [Validators.required]),
      alias: new FormControl(''),
      website: new FormControl('', /*[Validators.pattern('/[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi')]*/),
      contactInfo: this.getContactInfoFormGroup(),
      paymentTerm: new FormControl('', [Validators.required]),
      remark: new FormControl(''),
      gstn: new FormControl(''),
      status: new FormControl('', [Validators.required]),
      addresses: this.fb.array([
        this.getAddressFormGroup()
      ]),
      contactPersons: this.fb.array([
        this.getContactPersonFormGroup()
      ])
    })
    this.paymentTermsValues = Object.values(PaymentTerm);
    this.statusValues = Object.values(Status);
    this.addressTypeValues = Object.values(AddressType);
    this.personTypeValues = Object.values(PersonType);
  }

  getAddressControl(index: number): { [key: string]: AbstractControl } {
    return (this.addresses.controls[index] as FormGroup).controls
  }

  getContactPersonControl(index: number): { [key: string]: AbstractControl } {
    return (this.contactPersons.controls[index] as FormGroup).controls
  }
  getContactPersonContactInfoControl(index: number): { [key: string]: AbstractControl } {
    return ((this.contactPersons.controls[index] as FormGroup).controls['contactInfo'] as FormGroup).controls
  }
  getContactPersonAddressControl(index: number): { [key: string]: AbstractControl } {
    return ((this.contactPersons.controls[index] as FormGroup).controls['address'] as FormGroup).controls
  }

  get formControls(): { [key: string]: AbstractControl } {
    return this.updateVendorForm.controls
  }
  get companyContactInfo(): { [key: string]: AbstractControl } {
    return (this.updateVendorForm.controls['contactInfo'] as FormGroup).controls
  }
  get addresses(): FormArray {
    return this.updateVendorForm.controls['addresses'] as FormArray
  }
  addAddress() {
    this.addresses.push(this.getAddressFormGroup());
  }
  removeAddress(removeIndex: number) {
    this.addresses.removeAt(removeIndex);
  }

  get contactPersons(): FormArray {
    return this.updateVendorForm.controls['contactPersons'] as FormArray
  }
  addContactPerson() {
    this.contactPersons.push(this.getContactPersonFormGroup());
  }
  removeContactPerson(removeIndex: number) {
    this.contactPersons.removeAt(removeIndex);
  }
  showSpinner = true;
  vendorId: string = '';

  ngOnInit(): void {
    this.vendorId = this.route.snapshot.paramMap.get('vendorId') ?? ''

    this.vendorService.getVendorData(this.vendorId).subscribe(response => {
      if(response.data.addresses == undefined || response.data.addresses ==null){
        response.data.addresses = [];
      }
      if(response.data.contactPersons == undefined || response.data.contactPersons ==null){
        response.data.contactPersons = [];
      }
      if(response.data.contactPersons.length>1){
        for (let index = 1; index < response.data.contactPersons.length; index++){
          this.addContactPerson();
        }
      }

      if(response.data.addresses.length>1){
        for (let index = 1; index < response.data.contactPersons.length; index++){
          this.addAddress();
        }
      }
      this.updateVendorForm.patchValue(response.data);
      this.showSpinner = false;
    })
  }

  submitData() {
    this.showSpinner = true;
    this.vendorService.updateVendor(this.vendorId, this.updateVendorForm.value).subscribe({
      next: (data) => {

        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Vendor update.', 'Success', {
          disableTimeOut: false,
          timeOut:2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        console.log(`response from save `, data)
      }, error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Vendor was not updated.', 'Error', {
          disableTimeOut: false,
          timeOut:2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-danger alert-with-icon",
          positionClass: 'toast-top-right'
        });

        console.log(`error from save `, err)
      }, complete: () => {
        this.showSpinner = false;
        console.log(`save over`, this.updateVendorForm.value)
      }
    })
  }
}

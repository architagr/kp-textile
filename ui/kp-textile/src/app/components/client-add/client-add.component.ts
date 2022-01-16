import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ClientService } from 'src/app/services/client-service';

@Component({
  selector: 'app-client-add',
  templateUrl: './client-add.component.html',
  styleUrls: ['./client-add.component.scss']
})
export class ClientAddComponent implements OnInit {
  addClientForm: FormGroup;

  private getContactInfoFormGroup(): FormGroup {
    return this.fb.group({
      email: new FormControl(''),
      landline: new FormControl(''),
      mobile: new FormControl(''),
      whatsapp: new FormControl(''),
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
      mobile: new FormControl(''),
    });
  }

  private getContactPersonFormGroup(): FormGroup {
    return this.fb.group({
      salutation: new FormControl('', [Validators.required]),
      firstName: new FormControl('', [Validators.required]),
      lastName: new FormControl('', [Validators.required]),
      department: new FormControl('', [Validators.required]),
      personType: new FormControl('', [Validators.required]),
      remark: new FormControl('', [Validators.required]),
      address: this.getAddressFormGroup(),
      contactInfo: this.getContactInfoFormGroup(),
    })
  }
  constructor(
    private clientService: ClientService,
    private fb: FormBuilder
  ) {
    this.addClientForm = this.fb.group({
      companyName: new FormControl('', [Validators.required]),
      alias: new FormControl(''),
      website: new FormControl(''),
      contactInfo: this.getContactInfoFormGroup(),
      paymentTerm: new FormControl('', [Validators.required]),
      remark: new FormControl(''),
      status: new FormControl('', [Validators.required]),
      addresses: this.fb.array([
        this.getAddressFormGroup()
      ]),
      contactPersons: this.fb.array([
        this.getContactPersonFormGroup()
      ])
    })

  }
  get addresses(): FormArray {
    return this.addClientForm.controls['addresses'] as FormArray
  }
  addAddress() {
    this.addresses.push(this.getAddressFormGroup());
  }
  removeAddress(removeIndex: number) {
    this.addresses.removeAt(removeIndex);
  }

  get contactPersons(): FormArray {
    return this.addClientForm.controls['contactPersons'] as FormArray
  }
  addContactPerson() {
    this.contactPersons.push(this.getContactInfoFormGroup());
  }
  removeContactPerson(removeIndex: number) {
    this.contactPersons.removeAt(removeIndex);
  }

  ngOnInit(): void {

  }

  submitData(){
    console.log(this.addClientForm.value)
  }
}

import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
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
      mobile: new FormControl('',[Validators.pattern('^([0|\+[0-9]{1,5})?([7-9][0-9]{9})$')]),
    });
  }

  private getContactPersonFormGroup(): FormGroup {
    return this.fb.group({
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
    private clientService: ClientService,
    private fb: FormBuilder
  ) {
    this.addClientForm = this.fb.group({
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

  }

  getAddressControl(index: number):{ [key: string]: AbstractControl }{
    return (this.addresses.controls[index] as FormGroup).controls
  }

  getContactPersonControl(index: number):{ [key: string]: AbstractControl }{
    return (this.contactPersons.controls[index] as FormGroup).controls
  }
  getContactPersonContactInfoControl(index: number):{ [key: string]: AbstractControl }{
    return ((this.contactPersons.controls[index] as FormGroup).controls['contactInfo'] as FormGroup).controls
  }
  getContactPersonAddressControl(index: number):{ [key: string]: AbstractControl }{
    return ((this.contactPersons.controls[index] as FormGroup).controls['address'] as FormGroup).controls
  }

  get formControls(): { [key: string]: AbstractControl }{
    return this.addClientForm.controls
  }
  get companyContactInfo():{ [key: string]: AbstractControl }{
    return (this.addClientForm.controls['contactInfo']as FormGroup).controls
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
    this.clientService.addClient(this.addClientForm.value).subscribe(data=>{
      console.log(`response from save `, data)
    },err=>{
      console.log(`error from save `, err) 
    }, ()=>{
      console.log(`save over`, this.addClientForm.value) 
    })
  }
}

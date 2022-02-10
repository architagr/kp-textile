import { Component, Inject, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProductDto, QualityDto } from 'src/app/models/quality-model';

@Component({
  selector: 'app-quality-add',
  templateUrl: './quality-add.component.html',
  styleUrls: ['./quality-add.component.scss']
})
export class QualityAddComponent implements OnInit {
  addQualityForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<QualityAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {quality:QualityDto, products: ProductDto[]},
  ) {
    this.addQualityForm = new FormGroup({
      name:new FormControl('', [Validators.required]),
      id: new FormControl(''),
      hsnCode: new FormControl('', [Validators.required]),
      productId: new FormControl('', [Validators.required]),
    })
  }

  ngOnInit(): void {
    this.addQualityForm.patchValue(this.data.quality)
  }

  onNoClick(): void {
    this.dialogRef.close(null);
  }

  submit(): void{
    this.dialogRef.close(this.addQualityForm.value);
  }
  get formControls(): { [key: string]: AbstractControl } {
    return this.addQualityForm.controls
  }
}

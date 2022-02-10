import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProductDto } from 'src/app/models/quality-model';

@Component({
  selector: 'app-add-product',
  templateUrl: './add-product.component.html',
  styleUrls: ['./add-product.component.scss']
})
export class AddProductComponent implements OnInit {
  addProductForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<AddProductComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ProductDto,
  ) {
    this.addProductForm = new FormGroup({
      name:new FormControl(''),
      id: new FormControl('')
    })
  }

  ngOnInit(): void {
    this.addProductForm.patchValue(this.data)
  }

  onNoClick(): void {
    this.dialogRef.close(null);
  }
submit():void{
  this.dialogRef.close(this.addProductForm.value);

}
}

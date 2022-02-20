import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { GodownDto } from 'src/app/models/godown-model';

@Component({
  selector: 'app-add-godown',
  templateUrl: './add-godown.component.html',
  styleUrls: ['./add-godown.component.scss']
})
export class AddGodownComponent implements OnInit {

  addGodownForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<AddGodownComponent>,
    @Inject(MAT_DIALOG_DATA) public data: GodownDto,
  ) {
    this.addGodownForm = new FormGroup({
      name: new FormControl(''),
      id: new FormControl('')
    })
  }

  ngOnInit(): void {
    this.addGodownForm.patchValue(this.data)
  }

  onNoClick(): void {
    this.dialogRef.close(null);
  }
  submit(): void {
    this.dialogRef.close(this.addGodownForm.value);
  }
}

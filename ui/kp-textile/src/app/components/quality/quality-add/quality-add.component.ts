import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { QualityDto } from 'src/app/models/quality-model';

@Component({
  selector: 'app-quality-add',
  templateUrl: './quality-add.component.html',
  styleUrls: ['./quality-add.component.scss']
})
export class QualityAddComponent implements OnInit {
  addQualityForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<QualityAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: QualityDto,
  ) {
    this.addQualityForm = new FormGroup({
      name:new FormControl(''),
      id: new FormControl('')
    })
  }

  ngOnInit(): void {
    this.addQualityForm.patchValue(this.data)
  }

  onNoClick(): void {
    this.dialogRef.close(this.addQualityForm.dirty ? this.addQualityForm.value: null);
  }
}

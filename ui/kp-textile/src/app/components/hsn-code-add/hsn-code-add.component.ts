import { Component, OnInit, Inject } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';

@Component({
  selector: 'app-hsn-code-add',
  templateUrl: './hsn-code-add.component.html',
  styleUrls: ['./hsn-code-add.component.scss']
})
export class HsnCodeAddComponent implements OnInit {
  addHsnCodeForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<HsnCodeAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: HnsCodeDto,
  ) {
    this.addHsnCodeForm = new FormGroup({
      hsnCode:new FormControl(''),
      id: new FormControl('')
    })
  }

  ngOnInit(): void {
    this.addHsnCodeForm.patchValue(this.data)
  }

  onNoClick(): void {
    this.dialogRef.close(this.addHsnCodeForm.dirty ? this.addHsnCodeForm.value: null);
  }

}

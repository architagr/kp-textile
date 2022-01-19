import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ClientDto } from 'src/app/models/client-model';

@Component({
  selector: 'app-delete-confirmation',
  templateUrl: './delete-confirmation.component.html',
  styleUrls: ['./delete-confirmation.component.scss']
})
export class DeleteConfirmationComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<DeleteConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {heading: string, message: string},
  ) { }

  ngOnInit(): void {
  }

  onNoClick(): void {
    this.dialogRef.close(null);
  }
  onOkClick(): void {
    this.dialogRef.close(this.data);
  }

}

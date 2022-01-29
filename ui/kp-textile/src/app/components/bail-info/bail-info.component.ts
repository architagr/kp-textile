import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { BailInfoResponse } from 'src/app/models/item-model';

@Component({
  selector: 'app-bail-info',
  templateUrl: './bail-info.component.html',
  styleUrls: ['./bail-info.component.scss']
})
export class BailInfoComponent implements OnInit {
  displayedColumnsSales: string[] = ['BillNo', 'Date', 'BilledQuantity',];
  displayedColumns: string[] = [...this.displayedColumnsSales, 'ReceivedQuantity', 'RemainingQuantity'];

  constructor(
    public dialogRef: MatDialogRef<BailInfoComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { info: BailInfoResponse, baleName: string },
  ) { }

  ngOnInit(): void {
  }

}

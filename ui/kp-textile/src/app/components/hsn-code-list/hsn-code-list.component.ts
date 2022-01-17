import { Component, OnInit } from '@angular/core';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { HsnCodeService } from 'src/app/services/hsn-code-service';
import { MatDialog } from '@angular/material/dialog';
import { HsnCodeAddComponent } from '../hsn-code-add/hsn-code-add.component';
@Component({
  selector: 'app-hsn-code-list',
  templateUrl: './hsn-code-list.component.html',
  styleUrls: ['./hsn-code-list.component.scss']
})
export class HsnCodeListComponent implements OnInit {
  displayedColumns: string[] = ['HsnCode', 'Action'];
  hsnCodes: HnsCodeDto[] = []
  hsnCodesAll: HnsCodeDto[] = [];
  constructor(
    private hsnCodeService: HsnCodeService,
    public dialog: MatDialog
  ) { }

  ngOnInit(): void {
    this.getAllCodes();
  }

  getAllCodes() {
    this.hsnCodeService.getAllHsnCode().subscribe(data => {
      this.hsnCodes = data.data
      this.hsnCodesAll = data.data
    });
  }
  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    if (filterValue.length > 0)
      this.hsnCodes = this.hsnCodesAll.filter(x => x.hnsCode.toLocaleLowerCase().includes(filterValue.toLocaleLowerCase()));
    else
      this.hsnCodes = this.hsnCodesAll;
  }


  openDialog(hsnCode: HnsCodeDto): void {
    const dialogRef = this.dialog.open(HsnCodeAddComponent, {
      data: { id: "", hnsCode: "" } as HnsCodeDto,
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        const hsnCode = (result as HnsCodeDto)
        if (hsnCode.id === "") {
          this.addCode(hsnCode);
        }
      }
    });
  }

  addCode(hsnCode: HnsCodeDto) {
    this.hsnCodeService.addHsnCode(hsnCode).subscribe(response => {
      console.log(`data added `, response)
      this.getAllCodes();
    })
  }
}

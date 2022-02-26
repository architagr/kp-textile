import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { GodownDto, GodownListResponse, GodownResponse } from 'src/app/models/godown-model';
import { GodownService } from 'src/app/services/godown-service';
import { ToastService } from 'src/app/services/toast-service';
import { AddGodownComponent } from '../add-godown/add-godown.component';

@Component({
  selector: 'app-godown-list',
  templateUrl: './godown-list.component.html',
  styleUrls: ['./godown-list.component.scss']
})
export class GodownListComponent implements OnInit {
  displayedColumns: string[] = ['Name'];

  allGodown: GodownDto[] = [];
  constructor(private godownService: GodownService,
    public dialog: MatDialog,
    private toastService: ToastService
  ) { }

  ngOnInit(): void {
    this.getAllGodown(); 
  }

  getAllGodown(){
    this.allGodown = [];
    this.godownService.getAllGodown().subscribe({
      next: (response: GodownListResponse) => {
        this.allGodown = response.data
      }
    })
  }
  addGodownOpenDialog() {
    const dialogRef = this.dialog.open(AddGodownComponent, {
      data: { id: '', name: '' } as GodownDto,
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        const godown = (result as GodownDto)
        if (godown.id === "") {
          this.addGodown(godown);
        }
      }
    });
  }

  addGodown(godown: GodownDto) {
    this.godownService.addGodown(godown.name).subscribe({
      next:(response: GodownResponse)=>{
          console.log(`data added `, response)
          this.getAllGodown();
      },
      error: (err)=>{
        this.toastService.show("Error", err.error?.errorMessage ?? err.message)
      }
    })
  }

}

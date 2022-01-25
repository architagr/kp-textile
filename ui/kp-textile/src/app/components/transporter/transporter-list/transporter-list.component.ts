import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ToastrService } from 'ngx-toastr';
import { TransporterDto } from 'src/app/models/transporter-model';
import { TransporterService } from 'src/app/services/transporter-service';
import { DeleteConfirmationComponent } from '../../delete-confirmation/delete-confirmation.component';

@Component({
  selector: 'app-transporter-list',
  templateUrl: './transporter-list.component.html',
  styleUrls: ['./transporter-list.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*', })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
    trigger('detailExpand1', [
      state('collapsed', style({ 'padding-top': '0px', 'padding-bottom': '0px' })),
      state('expanded', style({ 'padding-top': '*', 'padding-bottom': '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class TransporterListComponent implements OnInit {
  filterForm: FormGroup
  searchText: string = ''
  displayedColumns: string[] = ['CompanyName', 'PaymentTerms', 'Status', 'ContactInfo', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 10;
  total: number = 0;
  lastEvalutionKey: any = null
  transporters: TransporterDto[] = []
  transporterAll: TransporterDto[] = []
  expandedElement: TransporterDto | null = null
  constructor(
    public transporterService: TransporterService,
    private toastr: ToastrService,
    public dialog: MatDialog
  ) { 
    this.filterForm = new FormGroup({
      searchText: new FormControl('')
    });

  }

  ngOnInit(): void {
    this.getTransporter();
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getTransporter();
  }

  get startIndex(): number{
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number{
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.transporterAll.length) {
      endIndex = this.transporterAll.length - 1
    }
    return endIndex
  }
  onNextPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;

    if (this.transporterAll.length >= this.endIndex + 1) {
      this.getTransporterFromLocalList()
    } else {
      this.getTransporter();
    }
  }
  getTransporterFromLocalList(){
    this.transporters = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex 
    for (; startIndex <= endIndex; startIndex++) {
      this.transporters.push(this.transporterAll[startIndex])
    }
  }
  onPrevPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getTransporterFromLocalList();
  }
  getTransporter() {
    this.transporterService.getAllTransporter(this.pageSize, this.searchText, this.lastEvalutionKey).subscribe(data => {
      this.transporters = data.data
      this.addToAllTransporterList(data.data);
      this.lastEvalutionKey = data.lastEvalutionKey
      this.total = data.total
    });
  }
  addToAllTransporterList(data: TransporterDto[]) {
    data.forEach(x => {
      if (!this.transporterAll.some(y => y.transporterId === x.transporterId)) {
        this.transporterAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getTransporter();
  }
  deleteTransporter(transporter: TransporterDto) {
    const dialogRef = this.dialog.open(DeleteConfirmationComponent, {
      data: { heading: 'Delete Transporter', message: `Are you sure you want to delete ${transporter.companyName} ?` },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        this.delete(transporter.transporterId);
      }
    });
  }
  delete(transporterId: string) {
    this.transporterService.deleteTransporter(transporterId).subscribe({
      next: (data) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Transporter deleted.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        this.getTransporter();
      },
      error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Error in deleting Transporter.', 'Error', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-danger alert-with-icon",
          positionClass: 'toast-top-right'
        });
      }
    })
  }

}

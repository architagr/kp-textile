import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { ClientDto } from 'src/app/models/client-model';
import { ClientService } from 'src/app/services/client-service';
import { ToastrService } from 'ngx-toastr';
import { MatDialog } from '@angular/material/dialog';
import { DeleteConfirmationComponent } from '../delete-confirmation/delete-confirmation.component';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-client-list',
  templateUrl: './client-list.component.html',
  styleUrls: ['./client-list.component.scss'],
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
export class ClientListComponent implements OnInit {
  filterForm: FormGroup
  searchText: string = ''
  displayedColumns: string[] = ['CompanyName', 'PaymentTerms', 'Status', 'ContactInfo', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 10;
  total: number = 0;
  lastEvalutionKey: any = null
  clients: ClientDto[] = []
  clientsAll: ClientDto[] = []
  expandedElement: ClientDto | null = null
  constructor(
    public clientService: ClientService,
    private toastr: ToastrService,
    public dialog: MatDialog
  ) {
    this.filterForm = new FormGroup({
      searchText: new FormControl('')
    });

  }

  ngOnInit(): void {
    this.getClients();
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getClients();
  }

  get startIndex(): number{
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number{
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.clientsAll.length) {
      endIndex = this.clientsAll.length - 1
    }
    return endIndex
  }
  onNextPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;

    if (this.clientsAll.length >= this.endIndex + 1) {
      this.getClientFromLocalList()
    } else {
      this.getClients();
    }
  }
  getClientFromLocalList(){
    this.clients = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex 
    for (; startIndex <= endIndex; startIndex++) {
      this.clients.push(this.clientsAll[startIndex])
    }
  }
  onPrevPageClick(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getClientFromLocalList();
  }
  getClients() {
    this.clientService.getAllClient(this.pageSize, this.searchText, this.lastEvalutionKey).subscribe(data => {
      this.clients = data.data
      this.addToAllClientList(data.data);
      this.lastEvalutionKey = data.lastEvalutionKey
      this.total = data.total
    });
  }
  addToAllClientList(data: ClientDto[]) {
    data.forEach(x => {
      if (!this.clientsAll.some(y => y.clientId === x.clientId)) {
        this.clientsAll.push(x);
      }
    });
  }
  applyFilter() {
    this.searchText = this.filterForm.controls['searchText'].value;
    this.lastEvalutionKey = null;
    this.getClients();
  }
  deleteClient(client: ClientDto) {
    const dialogRef = this.dialog.open(DeleteConfirmationComponent, {
      data: { heading: 'Delete Client', message: `Are you sure you want to delete ${client.companyName} ?` },
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed', result);
      if (result) {
        this.delete(client.clientId);
      }
    });
  }
  delete(clientId: string) {
    this.clientService.deleteClient(clientId).subscribe({
      next: (data) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Client deleted.', 'Success', {
          disableTimeOut: false,
          timeOut: 2000,
          closeButton: true,
          enableHtml: true,
          toastClass: "alert alert-success alert-with-icon",
          positionClass: 'toast-top-right'
        });
        this.getClients();
      },
      error: (err) => {
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Error in deleting client.', 'Error', {
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
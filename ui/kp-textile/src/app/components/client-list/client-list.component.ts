import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { ClientDto } from 'src/app/models/client-model';
import { ClientService } from 'src/app/services/client-service';
import { ToastrService } from 'ngx-toastr';
import { MatDialog } from '@angular/material/dialog';
import { DeleteConfirmationComponent } from '../delete-confirmation/delete-confirmation.component';

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
  displayedColumns: string[] = ['CompanyName', 'PaymentTerms', 'Status', 'ContactInfo', 'Action'];
  pageNumber: number = 0;
  pageSize: number = 1;
  total: number = 0;
  lastEvalutionKey: any = null
  clients: ClientDto[] = []
  clientsAll: ClientDto[] = []
  expandedElement: ClientDto | null = null
  constructor(
    public clientService: ClientService,
    private toastr: ToastrService,
    public dialog: MatDialog
  ) { }

  ngOnInit(): void {
    this.getClients();
  }
  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.lastEvalutionKey = null;
    this.pageNumber = 0;
    this.getClients();
  }
  onPageNumberChange(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getClients();
  }
  getClients() {
    this.clientService.getAllClient(this.pageSize, this.lastEvalutionKey).subscribe(data => {
      this.clients = data.data
      this.clientsAll = data.data
      this.lastEvalutionKey = data.lastEvalutionKey
      this.total = data.total
    });
  }
  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    if (filterValue.length > 0)
      this.clients = this.clientsAll.filter(x => x.companyName.toLocaleLowerCase().includes(filterValue.toLocaleLowerCase()));
    else
      this.clients = this.clientsAll;
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
        this.toastr.info('<span class="tim-icons icon-bell-55" [data-notify]="icon"></span> Clinet deleted.', 'Success', {
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
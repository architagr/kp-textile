import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { ClientDto } from 'src/app/models/client-model';
import { ClientService } from 'src/app/services/client-service';

@Component({
  selector: 'app-client-list',
  templateUrl: './client-list.component.html',
  styleUrls: ['./client-list.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*',})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
    trigger('detailExpand1', [
      state('collapsed', style({'padding-top': '0px', 'padding-bottom': '0px'})),
      state('expanded', style({'padding-top': '*', 'padding-bottom': '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class ClientListComponent implements OnInit {
  displayedColumns: string[] = ['CompanyName', 'Email', 'Address', 'ContactInfo', 'Action'];
  clients: ClientDto[] = []
  clientsAll: ClientDto[] = []
  expandedElement: ClientDto | null = null
  constructor(
    public clientService: ClientService
  ) { }

  ngOnInit(): void {
    this.clientService.getAllClient().subscribe(data => {
      this.clients = data.data
      this.clientsAll = data.data
    });
  }
  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    if(filterValue.length>0)
    this.clients = this.clients.filter(x=>x.companyName.toLocaleLowerCase().includes(filterValue.toLocaleLowerCase()))
    else
    this.clients = this.clientsAll
  }
}
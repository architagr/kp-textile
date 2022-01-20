import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { map, Observable, startWith } from 'rxjs';
import { ClientDto } from 'src/app/models/client-model';
import { HnsCodeDto } from 'src/app/models/hsn-code-model';
import { ClientService } from 'src/app/services/client-service';
import { HsnCodeService } from 'src/app/services/hsn-code-service';

@Component({
  selector: 'app-sales',
  templateUrl: './sales.component.html',
  styleUrls: ['./sales.component.scss']
})
export class SalesComponent implements OnInit {
  clients: ClientDto[] = [];
  hsnCodes: HnsCodeDto[] = [];
  salesForm: FormGroup;
  clinetFilteredOptions: Observable<ClientDto[]>;
  hsnCodeFilteredOptions: Observable<HnsCodeDto[]>;
  constructor(
    private hsnCodeService: HsnCodeService,
    private clientService: ClientService,
  ) {
    this.salesForm = new FormGroup({
      client: new FormControl(''),
      hsnCode: new FormControl(''),
    })
    this.clinetFilteredOptions = this.salesForm.controls['client'].valueChanges.pipe(
      startWith(''),
      map(value => this._filterClient(value)),
    );

    this.hsnCodeFilteredOptions = this.salesForm.controls['hsnCode'].valueChanges.pipe(
      startWith(''),
      map(value => this._filterHsnCode(value)),
    );

  }
  private _filterClient(clientSearchText: string): ClientDto[] {
    const filterValue = clientSearchText.toLowerCase();

    return this.clients.filter(option => option.companyName.toLowerCase().includes(filterValue));
  }

  private _filterHsnCode(hsnCodeSearchText: string): HnsCodeDto[] {
    const filterValue = hsnCodeSearchText.toLowerCase();

    return this.hsnCodes.filter(option => option.hnsCode.toLowerCase().includes(filterValue));
  }

  ngOnInit(): void {
    this.getClients();
    this.getHsnCodes();
  }
  getClients() {
    this.clientService.getAllClient(10,"", {}).subscribe(response => {
      this.clients = response.data
    });
    
  }
  getHsnCodes(){
    this.hsnCodeService.getAllHsnCode().subscribe(response => {
      this.hsnCodes = response.data
    });
  }
}

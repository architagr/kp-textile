import { Component, OnInit } from '@angular/core';
import { ROUTES } from 'src/app/models/route-model';

@Component({
  selector: 'app-menu',
  templateUrl: './menu.component.html',
  styleUrls: ['./menu.component.scss']
})
export class MenuComponent implements OnInit {
  menuItems: any[] = [];
  constructor() { }

  ngOnInit(): void {
    this.menuItems = ROUTES.filter(menuItem => menuItem.showMenu);
  }

}


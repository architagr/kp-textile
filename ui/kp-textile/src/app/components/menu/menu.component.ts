import { Component, OnInit } from '@angular/core';
import { ROUTES, RouteInfo } from 'src/app/models/route-model';
import { Location } from "@angular/common";
@Component({
  selector: 'app-menu',
  templateUrl: './menu.component.html',
  styleUrls: ['./menu.component.scss']
})
export class MenuComponent implements OnInit {
  menuItems: any[] = [];
  showMasterDataMenu:boolean = false;
  constructor(
    private location: Location,
  ) { }

  ngOnInit(): void {
    this.menuItems = ROUTES.filter(menuItem => menuItem.showMenu);
  }
  addActiveClass(route: RouteInfo): boolean {

    var titlee = this.location.prepareExternalUrl(this.location.path());
    if (titlee.charAt(0) === "#") {
      titlee = titlee.slice(1);   
    }
    
    if (titlee.indexOf(route.path) > -1 || (route.childRoutes.length>0 && route.childRoutes.some(x=>titlee.includes(x)))) {
      return true;
    }
    return false;
  }
}


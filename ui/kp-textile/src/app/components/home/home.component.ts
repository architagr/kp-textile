import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/services/spinner-service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  
isLogin: boolean = false;
  constructor(public spinnerService: SpinnerService,
    private router: Router) {
      if(router.url === '/' || router.url === '/login'){
        this.isLogin  = true
      }
    }
  // public sidebarColor: string = "red";
  // changeSidebarColor(color: string){
  //   var sidebar = document.getElementsByClassName('sidebar')[0];
  //   var mainPanel = document.getElementsByClassName('main-panel')[0];

  //   this.sidebarColor = color;

  //   if(sidebar != undefined){
  //       sidebar.setAttribute('data',color);
  //   }
  //   if(mainPanel != undefined){
  //       mainPanel.setAttribute('data',color);
  //   }
  // }
  // changeDashboardColor(color: string){
  //   var body = document.getElementsByTagName('body')[0];
  //   if (body && color === 'white-content') {
  //       body.classList.add(color);
  //   }
  //   else if(body.classList.contains('white-content')) {
  //     body.classList.remove('white-content');
  //   }
  // }
  ngOnInit(): void {
  }

}

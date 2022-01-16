import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientAddComponent } from './components/client-add/client-add.component';
import { ClientListComponent } from './components/client-list/client-list.component';
import { HomeComponent } from './components/home/home.component';

const routes: Routes = [
  {
    path:"",
    component:HomeComponent
  },
  {
    path:"dashboard",
    component:HomeComponent
  },
  {
    path:"client",
    component:ClientListComponent
  },
  {
    path:"addclient",
    component:ClientAddComponent
  },
  {
    path:"updateclient/:clientId",
    component:ClientAddComponent
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

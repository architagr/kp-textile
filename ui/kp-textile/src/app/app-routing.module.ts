import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientAddComponent } from './components/client-add/client-add.component';
import { ClientListComponent } from './components/client-list/client-list.component';
import { ClientUpdateComponent } from './components/client-update/client-update.component';
import { HomeComponent } from './components/home/home.component';
import { HsnCodeAddComponent } from './components/hsn-code-add/hsn-code-add.component';
import { HsnCodeListComponent } from './components/hsn-code-list/hsn-code-list.component';
import { SalesComponent } from './components/sales/sales.component';

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
    component:ClientUpdateComponent
  },
  {
    path:"hsncode",
    component: HsnCodeListComponent
  },
  {
    path:"addhsncode",
    component: HsnCodeAddComponent
  },
  {
    path:'sales',
    component: SalesComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

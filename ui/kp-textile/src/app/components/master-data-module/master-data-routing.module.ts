import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { GodownListComponent } from './components/godown/godown-list/godown-list.component';

const routes: Routes = [
    {
        path:'',
        component: GodownListComponent
    },
    {
      path:'godown',
      component: GodownListComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MasterDataRoutingModule { }

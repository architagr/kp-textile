import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { GodownListComponent } from './components/godown/godown-list/godown-list.component';
import { QualityListComponent } from './components/quality/quality-list/quality-list.component';

const routes: Routes = [

  {
    path: 'godown',
    component: GodownListComponent
  },
  {
    path: 'quality',
    component: QualityListComponent,
  },
  {
    path: '',
    redirectTo:'/godown',
    pathMatch:"full"
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MasterDataRoutingModule { }

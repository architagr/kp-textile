import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransporterListComponent } from './transporter-list.component';

describe('TransporterListComponent', () => {
  let component: TransporterListComponent;
  let fixture: ComponentFixture<TransporterListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TransporterListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TransporterListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

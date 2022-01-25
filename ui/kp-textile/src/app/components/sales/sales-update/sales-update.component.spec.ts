import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SalesUpdateComponent } from './sales-update.component';

describe('SalesUpdateComponent', () => {
  let component: SalesUpdateComponent;
  let fixture: ComponentFixture<SalesUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SalesUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SalesUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

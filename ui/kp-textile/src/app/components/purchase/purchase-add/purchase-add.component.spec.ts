import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PurchaseAddComponent } from './purchase-add.component';

describe('PurchaseAddComponent', () => {
  let component: PurchaseAddComponent;
  let fixture: ComponentFixture<PurchaseAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PurchaseAddComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PurchaseAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

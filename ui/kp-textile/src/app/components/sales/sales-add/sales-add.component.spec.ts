import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SalesAddComponent } from './sales-add.component';

describe('SalesAddComponent', () => {
  let component: SalesAddComponent;
  let fixture: ComponentFixture<SalesAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SalesAddComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SalesAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

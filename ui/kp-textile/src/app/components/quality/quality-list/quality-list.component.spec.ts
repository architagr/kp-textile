import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QualityListComponent } from './quality-list.component';

describe('QualityListComponent', () => {
  let component: QualityListComponent;
  let fixture: ComponentFixture<QualityListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QualityListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QualityListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

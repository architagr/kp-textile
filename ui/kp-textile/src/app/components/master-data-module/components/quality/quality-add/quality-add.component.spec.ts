import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QualityAddComponent } from './quality-add.component';

describe('QualityAddComponent', () => {
  let component: QualityAddComponent;
  let fixture: ComponentFixture<QualityAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QualityAddComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QualityAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

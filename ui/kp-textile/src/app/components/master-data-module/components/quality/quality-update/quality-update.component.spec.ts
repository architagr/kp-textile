import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QualityUpdateComponent } from './quality-update.component';

describe('QualityUpdateComponent', () => {
  let component: QualityUpdateComponent;
  let fixture: ComponentFixture<QualityUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QualityUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QualityUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

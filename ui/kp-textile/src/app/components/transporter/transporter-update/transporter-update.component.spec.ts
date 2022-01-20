import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransporterUpdateComponent } from './transporter-update.component';

describe('TransporterUpdateComponent', () => {
  let component: TransporterUpdateComponent;
  let fixture: ComponentFixture<TransporterUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TransporterUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TransporterUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

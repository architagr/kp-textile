import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BailInfoComponent } from './bail-info.component';

describe('BailInfoComponent', () => {
  let component: BailInfoComponent;
  let fixture: ComponentFixture<BailInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ BailInfoComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(BailInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

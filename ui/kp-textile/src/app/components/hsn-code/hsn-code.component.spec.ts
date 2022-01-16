import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HsnCodeComponent } from './hsn-code.component';

describe('HsnCodeComponent', () => {
  let component: HsnCodeComponent;
  let fixture: ComponentFixture<HsnCodeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HsnCodeComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HsnCodeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

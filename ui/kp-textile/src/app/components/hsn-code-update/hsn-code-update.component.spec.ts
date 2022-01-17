import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HsnCodeUpdateComponent } from './hsn-code-update.component';

describe('HsnCodeUpdateComponent', () => {
  let component: HsnCodeUpdateComponent;
  let fixture: ComponentFixture<HsnCodeUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HsnCodeUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HsnCodeUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

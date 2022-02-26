import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HsnCodeAddComponent } from './hsn-code-add.component';

describe('HsnCodeAddComponent', () => {
  let component: HsnCodeAddComponent;
  let fixture: ComponentFixture<HsnCodeAddComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HsnCodeAddComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HsnCodeAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddGodownComponent } from './add-godown.component';

describe('AddGodownComponent', () => {
  let component: AddGodownComponent;
  let fixture: ComponentFixture<AddGodownComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AddGodownComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AddGodownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GodownListComponent } from './godown-list.component';

describe('GodownListComponent', () => {
  let component: GodownListComponent;
  let fixture: ComponentFixture<GodownListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ GodownListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(GodownListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

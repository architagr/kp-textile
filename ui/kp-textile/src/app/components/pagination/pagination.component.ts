import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-pagination',
  templateUrl: './pagination.component.html',
  styleUrls: ['./pagination.component.scss']
})
export class PaginationComponent implements OnInit {

  form: FormGroup

  @Input('total') total: number = 0;

  @Input('pageSize') pageSize: number = 0;
  @Output('pageSizeChange') pageSizeChange: EventEmitter<number> = new EventEmitter<number>();

  @Input('pageNumber') pageNumber: number = 1;
  @Output('pageNumberChange') pageNumberChange: EventEmitter<number> = new EventEmitter<number>();

  constructor() {
    this.form = new FormGroup({
      pageSize: new FormControl(10)
    });
    this.form.get('pageSize')?.valueChanges.subscribe(val => {
      this.pageSize = parseInt(val);
      this.pageSizeChange.emit(this.pageSize);
    })
  }
  get start(): number {
    const startVal = (this.pageNumber * this.pageSize)+1
    return startVal
  }
  get end(): number {
    const endValue = (this.start + this.pageSize -1)
    if (endValue > this.total){
      return this.total
    }
    return endValue
  }
  ngOnInit(): void {
    this.form.patchValue({pageSize: this.pageSize})
  }
  onNextClick() {
    this.pageNumber++;
    this.pageNumberChange.emit(this.pageNumber);
  }
  onPreviousClick() {
    this.pageNumber--;
    this.pageNumberChange.emit(this.pageNumber);
  }
}

import { Component, OnInit } from '@angular/core';
import { BailDetailsDto, BailInfoResponse } from 'src/app/models/item-model';
import { QualityDto, QualityListItemDto } from 'src/app/models/quality-model';
import { BailService } from 'src/app/services/bail-service';
import { QualityService } from 'src/app/services/quality-serice';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { catchError, forkJoin, from, map, mergeMap, Observable, of, tap, toArray } from 'rxjs';

@Component({
  selector: 'app-quality-list',
  templateUrl: './quality-list.component.html',
  styleUrls: ['./quality-list.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*', })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
    trigger('detailExpand1', [
      state('collapsed', style({ 'padding-top': '0px', 'padding-bottom': '0px' })),
      state('expanded', style({ 'padding-top': '*', 'padding-bottom': '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class QualityListComponent implements OnInit {
  showSpinnerCount = 0;
  allQuality: QualityListItemDto[] = [];
  qualities: QualityListItemDto[] = [];
  expandedElement: QualityListItemDto | null = null
  displayedColumns: string[] = ['QualityName', 'RemainingQuantity', 'Action'];
  pageSize = 10;
  pageNumber = 0;

  constructor(
    private qualityService: QualityService,
    private bailService: BailService,
  ) { }

  ngOnInit(): void {
    this.getAllQuality();
  }
  getAllQuality() {
    this.showSpinnerCount++;
    this.qualityService.getAllQualities().subscribe({
      next: (data) => {
        this.showSpinnerCount--;
        this.getAllSalableBails(data.data);
      },
      complete: () => {

      }
    });
  }
  getAllSalableBails(qualities: QualityDto[]) {
    this.showSpinnerCount += qualities.length
    from(qualities).pipe(mergeMap((element) =>
      this.bailService.getBailInfoByQuality(element.id).pipe(catchError(error => {
        console.log(`error: `, { error })
        return of({} as BailInfoResponse)
      }),
      ),
      this.pageSize),
      toArray()).subscribe({
        next: (response) => {
          for (let index = 0; index < qualities.length; index++) {
            const element = qualities[index];
            const data = response.find(x=>x.statusCode===200 && x.purchase.length>0 && x.purchase[0].quality === element.id);

            if (data && data!.purchase!.length > 0) {
              let total = 0;
              data!.purchase.forEach(x => total = total + x.pendingQuantity);
              
              this.allQuality.push({ ...element, pendingQuantity: total, bailDetails: data!.purchase } as QualityListItemDto);
            } else {
              this.allQuality.push({ ...element, pendingQuantity: 0, bailDetails: [] } as QualityListItemDto);
            }
          }
          this.showSpinnerCount -=this.allQuality.length
          this.getQualitiesFromLocalList();
        }
      })
  
  }

  onPageSizeChange(pageSize: number) {
    this.pageSize = pageSize;
    this.pageNumber = 0;
  }

  get startIndex(): number {
    return this.pageNumber * this.pageSize
  }
  get endIndex(): number {
    let endIndex = this.startIndex + this.pageSize - 1
    if (endIndex > this.allQuality.length) {
      endIndex = this.allQuality.length - 1
    }
    return endIndex
  }
  onPageChange(pageNumber: number) {
    this.pageNumber = pageNumber;
    this.getQualitiesFromLocalList();
  }

  getQualitiesFromLocalList() {
    this.qualities = [];
    let startIndex = this.startIndex
    let endIndex = this.endIndex
    for (; startIndex <= endIndex; startIndex++) {
      this.qualities.push(this.allQuality[startIndex])
    }
  }
}

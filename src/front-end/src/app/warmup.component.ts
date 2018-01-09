import { Component, OnInit } from '@angular/core';
import { IDatapoint } from './datapoint';
import { EnergyService } from './energy.service';

@Component({
  selector: 'my-app',
  templateUrl: './warmup.component.html'
})
export class WarmupComponent {
  status: string = "Test";

  // currentDate: Date;
  // yesterday: Date;

  // datapoints: any[];

  // getLineArray(): number[] {
  //   return Array(15);
  // }

  // wattsToPixels(watts: number, pxOffset: number = 0): number {
  //   var px = watts / 40 + pxOffset;
  //   return px;
  // }
  initiateWarmup() {
    for (var i = 1; i < 10; i++) {
      this.status = i.toString();
      
    }
  }

  // previous(): void {
  //   const newDate = new Date(this.currentDate.valueOf());
  //   newDate.setDate(newDate.getDate() - 1)
  //   this.currentDate = newDate;
  //   this.updateDatapoints();
  // }

  // next(): void {
  //   const newDate = new Date(this.currentDate.valueOf());
  //   newDate.setDate(newDate.getDate() + 1)
  //   this.currentDate = newDate;
  //   this.updateDatapoints();
  // }

  // constructor(private _energyDataService: EnergyService) {
    
  // }

  // updateDatapoints(): void {
  //   this._energyDataService.getEnergyForDate(this.currentDate)
  //     .subscribe( datapoints => {
  //       this.datapoints = datapoints;
  //     },
  //     error => {}
  //   );
  // }

  // ngOnInit(): void {
  //   console.log('initing');
  //   this.currentDate = new Date();
  //   this.currentDate.setDate(this.currentDate.getDate() - 1);
  //   this.yesterday = this.currentDate;
  //   this.updateDatapoints();
  // }

}

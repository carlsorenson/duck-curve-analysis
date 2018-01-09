import { Component, OnInit } from '@angular/core';
import { IDatapoint } from './datapoint';
import { EnergyService } from './energy.service';

@Component({
  selector: 'my-app',
  templateUrl: './explore.component.html'
})
export class ExploreComponent implements OnInit {

  currentDate: Date;
  currentMonth: Date;

  yesterday: Date;
  startDate: Date;
  displayMode: string;
  datapoints: any[];

  displayModes = [
    {value: 'weekdays', text: 'Weekdays'},
    {value: 'weekends', text: 'Weekends'},
    {value: 'all', text: 'All days'},
    {value: 'single', text: 'Single days only'}
  ]

  onChangeDisplayMode(mode: string) {
    console.log(mode);
    this.displayMode = mode;
    this.updateDatapoints();
  }

  getLineArray(): number[] {
    return Array(15);
  }

  wattsToPixelsMonthly(watts: number, pxOffset: number = 0): number {
    var px = watts / 12 + pxOffset;
    return px;
  }

  wattsToPixelsSingle(watts: number, pxOffset: number = 0): number {
    var px = watts / 40 + pxOffset;
    return px;
  }

  previousMonth(): void {
    const newDate = new Date(this.currentMonth.valueOf());
    var month = newDate.getMonth();
    
    if (month == 1) {
      month = 12;
      const year = newDate.getFullYear() - 1;
      newDate.setFullYear(year);
    }
    else {
      month = month - 1;
    }

    newDate.setMonth(month)
    this.currentMonth = newDate;
    this.updateDatapoints();
  }

  nextMonth(): void {
    const newDate = new Date(this.currentMonth.valueOf());
    var month = newDate.getMonth();
    
    if (month == 12) {
      month = 1;
      const year = newDate.getFullYear() + 1;
      newDate.setFullYear(year);
    }
    else {
      month = month + 1;
    }

    newDate.setMonth(month)
    this.currentMonth = newDate;
    this.updateDatapoints();
  }

  previousDay(): void {
    const newDate = new Date(this.currentDate.valueOf());
    newDate.setDate(newDate.getDate() - 1)
    this.currentDate = newDate;
    this.updateDatapoints();
  }

  nextDay(): void {
    const newDate = new Date(this.currentDate.valueOf());
    newDate.setDate(newDate.getDate() + 1)
    this.currentDate = newDate;
    this.updateDatapoints();
  }

  constructor(private _energyDataService: EnergyService) {
    
  }

  updateDatapoints(): void {
    switch (this.displayMode) {
    case 'weekdays':
    case 'weekends':
    case 'all':
        this._energyDataService.getEnergyAverages(this.displayMode, this.currentMonth)
        .subscribe( datapoints => {
          this.datapoints = datapoints;
        },
        error => {}
      );
      break;
    case 'single':
        this._energyDataService.getEnergyForDate(this.currentDate)
        .subscribe( datapoints => {
          this.datapoints = datapoints;
        },
        error => {}
      );
      break;
    }
    
  }

  ngOnInit(): void {
    console.log('initing');
    this.currentDate = new Date();
    this.currentDate.setHours(0);
    this.currentDate.setMinutes(0);
    this.currentDate.setSeconds(0);
    this.currentDate.setDate(this.currentDate.getDate() - 1);

    this.yesterday = this.currentDate;

    this.currentMonth = new Date(this.currentDate.getTime());
    this.currentMonth.setDate(1);

    this.startDate = new Date(2017, 8, 15);

    this.displayMode = 'weekdays';
    this.updateDatapoints();
  }

}

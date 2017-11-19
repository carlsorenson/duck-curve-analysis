import { Component } from '@angular/core';
import { IDatapoint } from './datapoint';
import { EnergyService } from './energy.service';

@Component({
  selector: 'my-app',
  templateUrl: './app.component.html'
})
export class AppComponent  { 
  
  currentDate: Date;
  yesterday: Date;
  datapoints: IDatapoint[];
  
  previous(): void {
    console.log('in previous');
    this.currentDate.setDate(this.currentDate.getDate() - 1);
    this.currentDate = this.currentDate;
    this.updateDatapoints();
  }

  constructor(private _energyDataService: EnergyService) {
    
  }

  updateDatapoints(): void {
    this._energyDataService.getEnergyForDate(this.currentDate)
      .subscribe( datapoints => {
        this.datapoints = datapoints;
      },
      error => {}
    );
  }

  ngOnInit(): void {
    console.log('initing');
    this.currentDate = new Date();
    this.currentDate.setDate(this.currentDate.getDate() - 1);
    this.yesterday = this.currentDate;
    this.updateDatapoints();
  }

}

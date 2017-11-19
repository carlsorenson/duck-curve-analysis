import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/map';

import { IDatapoint } from './datapoint';

@Injectable()
export class EnergyService {
  private _apiUrl = 'https://duck-curve-analysis.appspot.com/api/v1/day/';

  constructor(private _http: HttpClient) { }

  zeroPad(value: number) {
    let stringValue = value.toString();
    if (stringValue.length < 2) {
        stringValue = "0" + stringValue
    }
    return stringValue
  }

  getEnergyForDate(selectedDate: Date): Observable<IDatapoint[]> {
    const url = this._apiUrl + this.zeroPad(selectedDate.getUTCFullYear()) + "-" + this.zeroPad(selectedDate.getUTCMonth() + 1) + "-" + this.zeroPad(selectedDate.getUTCDate())
    console.log(selectedDate, selectedDate.getDate())
    console.log('hi!', url);
    return this._http.get<IDatapoint[]>(url)
        .do(data => console.log('All: ' + JSON.stringify(data)))
        .catch(this.handleError);
  }

  private handleError(err: HttpErrorResponse) {
      // in a real world app, we may send the server to some remote logging infrastructure
      // instead of just logging it to the console
      let errorMessage = '';
      if (err.error instanceof Error) {
          // A client-side or network error occurred. Handle it accordingly.
          errorMessage = `An error occurred: ${err.error.message}`;
      } else {
          // The backend returned an unsuccessful response code.
          // The response body may contain clues as to what went wrong,
          errorMessage = `Server returned code: ${err.status}, error message is: ${err.message}`;
      }
      console.error(errorMessage);
      return Observable.throw(errorMessage);
  }

}

import { NgModule }      from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';

import { EnergyService } from './energy.service'

import { AppComponent }  from './app.component';
import { HomeComponent }  from './home.component';
import { ExploreComponent }  from './explore.component';
import { ConclusionsComponent }  from './conclusions.component';
import { WarmupComponent } from './warmup.component';

@NgModule({
  imports:      [ 
    BrowserModule, 
    HttpClientModule,
    RouterModule.forRoot([
      { path: 'explore', component: ExploreComponent },
      { path: 'conclusions', component: ConclusionsComponent },
      { path: 'warmup', component: WarmupComponent },
      { path: '', component: HomeComponent }
    ])
  ],
  declarations: [ AppComponent, HomeComponent, ExploreComponent, ConclusionsComponent, WarmupComponent ],
  bootstrap:    [ AppComponent ],
  providers: [
    EnergyService    
  ]
})
export class AppModule { }

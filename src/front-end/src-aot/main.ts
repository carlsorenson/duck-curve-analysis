import { platformBrowser }    from '@angular/platform-browser';
import { AppModuleNgFactory } from '../aot/intermediate/app/app.module.ngfactory';

console.log('Running AOT compiled');
platformBrowser().bootstrapModuleFactory(AppModuleNgFactory);

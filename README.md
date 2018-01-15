# duck-curve-analysis
### Here be dragons!
I used this project to give me experience with Go and Angular. It's still a work in progress, so there are many things that should be fixed and cleaned up. However, I am pleased to present this duck curve analysis tool. See the [Home page](https://duck-curve-analysis.appspot.com/) of the site for a discussion of what the duck curve is. Then, click to the other pages to see the tool in action and to read about the results and the conclusions we can draw. 

### Technologies used
The codebase is a combination of a Go back end and an Angular front end. I managed to set it up so that I can develop the Angular in JIT-compiled mode locally, and then run an AOT compilation for deployment. The deployment script gathers up the files that will be needed and deploys both the AOT-compiled front end and the Go back end to Google AppEngine.

The data is displayed using some embedded SVG. I found SVG and Angular to be a fantastic fit - as the data is updated, the chart changes along with it. 

The back end is a simple Go API. It abstracts the calls to the Neurio API, so that the front end can make simpler calls. It caches the results in Google Datastore so that we can quickly get all the data needed to produce monthly averages. The back end does the calculations for the averages, keeping the front end focused on simply displaying the data and options to the user.

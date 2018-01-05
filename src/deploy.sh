#!/bin/bash

gcloud config set project duck-curve-analysis
cp -r front-end/src/. front-end/intermediate/
rm front-end/intermediate/index-jit.html
cp -r front-end/src-aot/. front-end/intermediate/
npm run build:aot --prefix front-end
gcloud app deploy --quiet

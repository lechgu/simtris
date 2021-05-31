#!/bin/sh

cd front
npm install
node_modules/webpack/bin/webpack.js --mode production

#!/bin/sh
export GOOS=js
export GOARCH=wasm
fswatch -0 -or . | xargs -0 -n 1 -I {} 'go' 'build' '-o' 'main.wasm'

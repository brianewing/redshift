#!/bin/sh

set -o errexit
set -o xtrace

export GOOS=js
export GOARCH=wasm

go build -o redshift.wasm

echo "Built redshift.wasm successfully."
echo "You will need to copy it and the $(go env GOROOT)/misc/wasm/wasm_exec.js file in order to use it in your project."

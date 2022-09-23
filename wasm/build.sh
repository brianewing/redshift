#!/bin/sh

set -e

OUTFILE=${OUTFILE:-ledplane.wasm}

GOOS=js GOARCH=wasm go build -v -o $OUTFILE $@

cp ../agpl-3.0.txt $OUTFILE.license.txt

echo "Build success!\n"

echo "Copy $OUTFILE, $OUTFILE.license.txt and $(go env GOROOT)/misc/wasm/wasm_exec.js to use Ledplane in your project."


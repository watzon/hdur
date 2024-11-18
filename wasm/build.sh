#!/bin/bash

# Create dist directory if it doesn't exist
mkdir -p dist

# Copy wasm_exec.js from Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

# Build the wasm binary with size optimizations
GOOS=js GOARCH=wasm go build \
    -ldflags="-s -w" \
    -trimpath \
    -o dist/hdur.wasm main.go

# Use wasm-opt for additional optimization if available
if command -v wasm-opt &> /dev/null; then
    echo "Optimizing with wasm-opt..."
    wasm-opt -Oz --enable-bulk-memory -o dist/hdur.wasm.opt dist/hdur.wasm
    mv dist/hdur.wasm.opt dist/hdur.wasm
else
    echo "wasm-opt not found, skipping additional optimization"
fi

# Copy our JavaScript wrapper and types
cp hdur.js dist/
cp hdur.d.ts dist/

# Show size comparison
echo "WASM file size:"
ls -lh dist/hdur.wasm

echo "Build complete! Files are in the dist directory."

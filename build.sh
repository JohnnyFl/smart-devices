#!/bin/bash

set -e

echo "Building Lambda functions..."

# Function names
FUNCTIONS=("get-device" "list-devices" "create-device" "update-device" "delete-device" "sqs-listener")

# Clean previous builds
rm -rf build
mkdir -p build

# Build each Lambda function
for fn in "${FUNCTIONS[@]}"; do
  echo "Building $fn..."

  # Build Go binary to bin/<fn>/bootstrap
  mkdir -p bin/$fn
  GOOS=linux GOARCH=amd64 go build -o bin/$fn/bootstrap cmd/$fn/main.go

  # Prepare clean zip structure
  mkdir -p build/$fn
  cp bin/$fn/bootstrap build/$fn/
  chmod +x build/$fn/bootstrap

  # Zip the bootstrap so it's at the root of the ZIP
  (cd build/$fn && zip -q ../$fn.zip bootstrap)
done

echo "All functions built and zipped successfully!"

#!/bin/bash

echo "Building Cliper..."
go build -o Cliper ./cmd/cliper

if [ $? -eq 0 ]; then
    echo "Build successful! Application created as 'Cliper'"
    echo "To run: ./Cliper"
else
    echo "Build failed."
fi
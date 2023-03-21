#!/bin/bash

go build -o main
echo "Build Done......."

# shellcheck disable=SC2028
echo "Output:"
./main

echo "Done..."
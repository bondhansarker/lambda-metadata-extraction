#!/bin/bash

go build -o main
echo "Build Done......."

echo "Zipping the binary...."
zip main.zip main

echo "Done..."
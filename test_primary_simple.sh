#!/bin/bash

cd /home/mleku/src/github.com/mleku/gio/cmd/primary || exit 1

echo "Building primary demo..."
go build -o primary-demo . 

if [ $? -eq 0 ]; then
    echo "Build successful. Running demo..."
    echo ""
    echo "Test steps:"
    echo "1. Type 'Hello World'"
    echo "2. Select 'Hello' with mouse"
    echo "3. Middle-click elsewhere to paste"
    echo ""
    ./primary-demo
else
    echo "Build failed"
    exit 1
fi



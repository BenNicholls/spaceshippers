#!/bin/bash

echo "Cleaning up!"

go fmt 

echo "Building Spaceshippers for Release"
echo "Remember to turn off DEBUG you clod"

go install -ldflags '-s -w -H=windowsgui'

echo "Yay, release is built! Maybe someday later make this script package everything up for you???"
#!/bin/bash

echo "Cleaning up!"

go fmt 

echo "Building Spaceshippers"

go install && spaceshippers

echo "Yay, we had fun!"
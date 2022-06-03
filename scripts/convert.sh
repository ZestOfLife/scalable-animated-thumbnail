#!/bin/bash

# Set frames to 4:3 aspect ratio

eval "mogrify frames-$2/*.jpg -gravity center -crop 4:3 frames-$2/*.jpg"

# Convert into gif and move to host path

eval "convert frames-$2/*.jpg videos/$1"

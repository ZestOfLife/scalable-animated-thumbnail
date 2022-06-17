#!/bin/bash

# Get frame folder
frame_name=$(basename -- "$1")
frame_name="${frame_name%.*}"

# Convert into gif and move to host path

eval "convert videos/$2 videos/frames-$frame_name/$3 videos/$2"

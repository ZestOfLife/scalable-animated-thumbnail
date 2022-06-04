#!/bin/bash

# Get frame folder
frame_name=$(basename -- "$1")
frame_name="${frame_name%.*}"

# Get frames
eval "scripts/frames.sh $1 $2"

# Resize frames
eval "scripts/resize.sh $1 *.jpg"

# Convert frames
eval "scripts/convert.sh $1 $2 *.jpg"

# Cleanup
eval "rm -r videos/frames-$frame_name"

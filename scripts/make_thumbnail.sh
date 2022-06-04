#!/bin/bash

# Make dir to store frames
frame_name=$(basename -- "$1")
frame_name="${frame_name%.*}"
eval "mkdir frames-${frame_name}"

# Get frames
eval "scripts/frames.sh $1 $frame_name"

# Convert frames
eval "scripts/convert.sh $2 $frame_name"

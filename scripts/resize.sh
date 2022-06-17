#!/bin/bash

# Get frame folder
frame_name=$(basename -- "$1")
frame_name="${frame_name%.*}"

# Set frames to 4:3 aspect ratio

eval "mogrify videos/frames-$frame_name/$2 -resize 720x540 videos/frames-$frame_name/$2"

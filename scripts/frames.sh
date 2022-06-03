#!/bin/bash

# Pull duration and get 2/3 of the frames
duration=$(eval "ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 videos/$1")
duration=$(echo $duration | cut -d "." -f 1 | cut -d "," -f 1)

dist_from=$((duration*2/3))
dist_to=$((dist_from+15)) # Around 15 seconds worth of frames we collect

eval "ffmpeg -i videos/$1 -vf \"select=between(t\\,$dist_from\\,$dist_to)\" -r 1 \"frames-$2/0%03d.jpg\"" # We do not want to build a mini video, rather we take a few frames for each second

#!/bin/bash
docker run --rm -it \
-e MTX_RTSPTRANSPORTS=tcp \
-e MTX_WEBRTCADDITIONALHOSTS=192.168.73.130 \
-v "$(pwd)/mediamtx.yml:/mediamtx.yml" \
-p 8554:8554 \
-p 1935:1935 \
-p 8888:8888 \
-p 8889:8889 \
-p 8890:8890/udp \
-p 8189:8189/udp \
--network host \
--name mediamtx \
bluenviron/mediamtx:latest-ffmpeg
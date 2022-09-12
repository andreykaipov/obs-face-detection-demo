#!/bin/sh

export OBS_FACE_SOURCE=Webcam
export OBS_FACE_CHECK_THRESHOLD=3
export OBS_FACE_CHECK_INTERVAL=5
export OBS_SCENE_BRB='zzz BRB'
export OBS_HOST="$WSL_HOST:4444"
export OBS_PASSWORD=hello

go run ./...

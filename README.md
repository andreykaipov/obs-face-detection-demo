## what is this?

It's a quick demo of automatic scene swapping in OBS using face detection.

For example, suppose you are streaming but you have to step away from the
computer because you forgot to put away the milk carton and now your cat has
jumped up on the kitchen counter top, knocking it over and making an absolute
mess. You don't have time to manually swap to a new scene to notify your
audience that you've stepped away and will be back shortly, so instead you use
this tool to do it for you!

## usage

```console
❯ OBS_CHECK_INTERVAL=2 \
  OBS_FACE_SOURCE=Webcam \
  OBS_SCENE_BRB='zzz BRB' \
  OBS_HOST="$WSL_HOST:4444" \
  OBS_PASSWORD=hello \
  go run ./...
```

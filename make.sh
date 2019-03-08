#!/bin/sh

export APP_IMAGE=jeremiaswerner/helloworld
export REBUILD_IMAGE=jeremiaswerner/rebuild

go build -o /dev/null helloworld.go  # Fail fast for compilation errors
docker build -t $APP_IMAGE .       # Do the real build and create image
docker push $APP_IMAGE

go build -o /dev/null rebuild.go  # Fail fast for compilation errors
docker build -t $REBUILD_IMAGE -f Dockerfile.rebuild .     # Do the real build and create image
docker push $REBUILD_IMAGE

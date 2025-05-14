#!/bin/bash

set -e

cd $CONTEXT_DIR
rm /tmp/build_args || echo OK
env >/tmp/build_args
echo "--build-arg \""$(cat /tmp/build_args | sed -z 's/\n/" --build-arg "/g')"IGNORE_VAR=IGNORE_VAR\"" >/tmp/build_args
BUILD_ARGS=$(cat /tmp/build_args)
echo "ADDITIONAL_BUILD_ARGS: $ADDITIONAL_BUILD_ARGS "
COMMAND="docker build -t $FULL_IMAGE_NAME $ADDITIONAL_BUILD_ARGS -f ./Dockerfile $BUILD_ARGS --no-cache ."
/bin/bash -c "$COMMAND"
docker push $FULL_IMAGE_NAME
rm /tmp/build_args
#!/bin/bash
set -e -x
npm -v
if [ $? -eq 0 ]
then
    cd $1/taskman-ui
    npm install --legacy-peer-deps
    npm run plugin
else
    docker run --rm -v $1:/app/taskman --name wetaskman-node-build node:12.13.1 /bin/bash /app/taskman/build/build-ui-docker.sh
fi
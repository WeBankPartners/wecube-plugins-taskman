#!/bin/bash
set -e -x
npm -v
if [ $? -eq 0 ]
then
    cd $1/taskman-ui
    npm install
    npm run build
    cd dist
    mkdir -p wetaskman
    mv js css img fonts favicon.ico taskman/
    cd ..
    mv dist dist_tmp
    npm run plugin
    mv dist plugin
    mv dist_tmp dist
else
    docker run --rm -v $1:/app/taskman --name wetaskman-node-build node:12.13.1 /bin/bash /app/taskman/build/build-ui-docker.sh
fi
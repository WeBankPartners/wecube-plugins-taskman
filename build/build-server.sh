#!/bin/bash
set -e -x
cd $(dirname $0)/../taskman-server
go build -ldflags "-linkmode external -extldflags -static -s"
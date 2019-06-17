#!/usr/bin/env bash

deps(){
# list all deps
go list -f '{{ join .Deps "\n" }}' | tee deps.log
}

ARCH=amd64
ARCH=x86
OS=darwin
LDFLAGS="-s -w "
OPTS=""


# build the plugins

pushd impl/apps/demo >/dev/null
GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build  -tags=gocql_debug -ldflags "$LDFLAGS" -buildmode=plugin -o demo-app-plugin.so github.com/hedzr/voxr-lite/misc/impl/apps/demo/... && ls -aGl *.so && file *.so
popd  >/dev/null

pushd impl/filters/pre >/dev/null
GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build  -tags=gocql_debug -ldflags "$LDFLAGS" -buildmode=plugin -o default-pre-filter.so github.com/hedzr/voxr-lite/misc/impl/filters/pre/... && ls -aGl *.so && file *.so
popd  >/dev/null

pushd impl/filters/post >/dev/null
GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build  -tags=gocql_debug -ldflags "$LDFLAGS" -buildmode=plugin -o default-post-filter.so github.com/hedzr/voxr-lite/misc/impl/filters/post/... && ls -aGl *.so && file *.so
popd  >/dev/null

#
# file:///Users/hz/hzw/golang-dev/src/github.com/hedzr/voxr-lite/misc/impl/filters/pre/default-pre-filter.so
# file:///Users/hz/hzw/golang-dev/src/github.com/hedzr/voxr-lite/misc/impl/filters/post/default-post-filter.so
#




#!/usr/bin/env bash


# list all deps
go list -f '{{ join .Deps "\n" }}' | tee deps.log


# build the plugin
cd impl/apps/demo
go build -buildmode=plugin -o demo-plugin.so demo.go








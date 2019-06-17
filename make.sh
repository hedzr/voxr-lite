#!/usr/bin/env bash

build-core() {
	PKG_SRC=./core/cli/main.go APPNAME=vx-core ./build.sh $*
}

build-misc() {
	PKG_SRC=./misc/cli/main.go APPNAME=vx-misc ./build.sh $*
}

build-circle() {
	PKG_SRC=./circle/cli/main.go APPNAME=vx-circle ./build.sh $*
}

build-all() {
	for n in core misc circle; do
		APPNAME=vx-$n PKG_SRC=./$n/cli/main.go ./build.sh all $*
	done
}

build-full() {
	for n in core misc circle; do
		APPNAME=vx-$n PKG_SRC=./$n/cli/main.go ./build.sh full $*
	done
}

build-all-linux() {
	for n in core misc circle; do
		APPNAME=vx-$n PKG_SRC=./$n/cli/main.go ./build.sh linux $*
	done
}

build-ci() {
  go mod download
	for n in core misc circle; do
		APPNAME=vx-$n PKG_SRC=./$n/cli/main.go ./build.sh all $*
	done

	ls -la ./bin/
	for f in bin/*; do gzip $f; done 
	ls -la ./bin/
}




mysql(){
	local PWD='Cofy#A99Izol'
	mysql --protocol=tcp -P 33306 -uroot -p'Cofy#A99Izol' im_core
}




build-extras() {
	build-apps
	build-filters
	true
}

build-apps() {
	ARCH=amd64
	# ARCH=x86
	#OS=darwin
	#OS=linux
	
	LDFLAGS="-s -w "
	OPTS=""
	
	pushd misc/impl/apps/demo >/dev/null
	rm *.so; GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go version; # GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 \
	go build -ldflags "$LDFLAGS" -buildmode=plugin    -o demo-app-plugin.so ./...  && ls -aGl *.so && file *.so
	# pushd misc/impl/apps/demo; go build -ldflags "-s -w" -buildmode=plugin -o demo-app-plugin.so ./...; popd
	popd >/dev/null
}

build-filters() {
	ARCH=amd64
	# ARCH=x86
	#OS=darwin
	#OS=linux

	LDFLAGS="-s -w "
	OPTS=""
	for t in pre post; do
	  d=misc/impl/filters/$t
		pushd $d >/dev/null
		rm *.so; # GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 \
		go build -tags=gocql_debug -ldflags "$LDFLAGS" -buildmode=plugin    -o demo-$t-filter-plugin.so ./...  && ls -aGl *.so && file *.so
		popd >/dev/null
	done
}



run-core() {
	go run ./core/cli/main.go $*
}

run-misc() {
	go run ./misc/cli/main.go $*
}



fmt() {
	gofmt -l -w -s .
}

lint() {
  golint ./...
}

gotest() {
  go test ./...
}

test() {
  go test ./...
}

gocov() {
  go test -race -covermode=atomic -coverprofile cover.out && \
  go tool cover -html=cover.out -o cover.html && \
  open cover.html
}

gocov-codecov() {
  # https://codecov.io/gh/hedzr/cmdr
  go test -race -coverprofile=coverage.txt -covermode=atomic
  bash <(curl -s https://codecov.io/bash) -t $CODECOV_TOKEN
}

gocov-codecov-open() {
  open https://codecov.io/gh/hedzr/cmdr
}


[[ $# -eq 0 ]] && {
	run-core
} || {
	cmd=$1 && shift
	case $cmd in
	*) $cmd "$@" ;;
	esac
}


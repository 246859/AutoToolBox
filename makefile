# basic info
app := tbm
module := github.com/246859/AutoToolBox/v3/cmd/toolboxmenu
output := $(shell pwd)/build
# meta info
git_version := $(shell git tag --sort=-version:refname | sed -n 1p)
# build info
host_os := $(shell go env GOHOSTOS)
host_arch := $(shell go env GOHOSTARCH)
os := $(host_os)
arch := $(host_arch)

ifeq ($(os), windows)
	exe := .exe
endif


.PHONY: build
build:
	# go lint
	go vet ./...

	# prepare target environment $(os)/$(arch)
	go env -w GOOS=$(os)
	go env -w GOARCH=$(arch)

	# build go module
	go build -trimpath \
		-ldflags="-X main.AppName=$(app) -X main.Version=$(git_version)" \
		-o $(output)/$(app)-$(os)-$(arch)$(exe) \
		$(module)

	# resume host environment $(host_os)/$(host_arch)
	go env -w GOOS=$(host_os)
	go env -w GOARCH=$(host_arch)
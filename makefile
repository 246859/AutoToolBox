app_name := toolbox
app_package := github.com/246859/AutoToolBox/v2/cmd/toolbox
hostos := $(shell go env GOHOSTOS)
hostarch := $(shell go env GOHOSTARCH)
os =
arch =
exe =

ifeq ($(os),)
	os := $(shell go env GOOS)
endif
ifeq ($(arch),)
	arch := $(shell go env GOARCH)
endif
ifeq ($(os),windows)
	exe = .exe
endif

bin := $(app_name)-$(os)-$(arch)$(exe)

.PHONY: build
build:
	# set target environment
	go env -w GOOS=$(os)
	go env -w GOARCH=$(arch)

	# build binary file
	go vet ./...
	go build -trimpath -o ./build/$(bin) $(app_package)

	# resume host environment
	go env -w GOOS=$(hostos)
	go env -w GOARCH=$(hostarch)

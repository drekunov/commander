VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "undefined")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "undefined")
BRANCH  ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "undefined")
BUILD_TIME ?= $(shell date -u +%s)

PKG := github.com/drekunov/gc/internal/app

LD_FLAGS := -ldflags "\
	-X '$(PKG).version=$(VERSION)' \
	-X '$(PKG).commit=$(COMMIT)' \
	-X '$(PKG).branch=$(BRANCH)' \
	-X '$(PKG).buildUnixTimestamp=$(BUILD_TIME)'"

lint:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	# golangci-lint 2.12.2's staticcheck cannot parse the Go 1.27 export data
	# format; pin a toolchain it supports.
	GOTOOLCHAIN=go1.26.4 golangci-lint run

modup:
	go get -u ./...
	go mod tidy

run:
	go run ./cmd/commander/main.go

pvs:
	pvs-golang analyze .

build:
	go build $(LD_FLAGS) -o ./commander ./cmd/commander/main.go
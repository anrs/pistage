NS := github.com/projecteru2/pistage
BUILD := go build -race
TEST := go test -count=1 -race -cover

LDFLAGS += -X "$(NS)/ver.Git=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(NS)/ver.Compile=$(shell go version)"
LDFLAGS += -X "$(NS)/ver.Date=$(shell date +'%F %T %z')"

PKGS := $$(go list ./...)

.PHONY: all test build grpc

default: build

rundev: build
	bin/pistaged --config dev.toml server

build: build-srv build-cli

build-srv:
	$(BUILD) -ldflags '$(LDFLAGS)' -o bin/pistaged cmd/server/server.go

build-cli:
	$(BUILD) -ldflags '$(LDFLAGS)' -o bin/pistagec cmd/client/client.go

lint: fmt
	golint $(PKGS)
	golangci-lint run

fmt: vet
	gofmt -s -w $$(find . -iname '*.go' | grep -v -P '\./vendor/')

vet:
	go vet $(PKGS)

deps:
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

test:
ifdef RUN
	$(TEST) -v -run='${RUN}' $(PKGS)
else
	$(TEST) $(PKGS)
endif

grpc:
	protoc --go_out=plugins=grpc:. grpc/gen/pistaged.proto

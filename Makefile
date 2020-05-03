PROTOS := $(wildcard *.proto) $(wildcard */*.proto) $(wildcard */*/*.proto)
PBGO := $(PROTOS:.proto=.pb.go)

all: $(PBGO)
	go build .

include .make/golangci-lint.mk
include .make/protoc.mk
include .make/grpc-go.mk

tidy:
	go mod tidy

lint: $(GOLANGCI_LINT)
	$(realpath $(GOLANGCI_LINT)) run

clean/proto:
	rm -f $(PBGO)

clean: clean/golangci-lint clean/protoc clean/protoc-gen-go clean/proto
	rm -f go.sum

test: # -count=1 disables cache
	go test -v -race -count=1 .

.PHONY: all tidy lint clean test

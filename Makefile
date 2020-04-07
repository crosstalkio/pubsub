PROTOS := $(wildcard *.proto) $(wildcard */*.proto)
PBGO := $(PROTOS:.proto=.pb.go)

all: $(PBGO)
	go build .

tidy:
	go mod tidy

clean: clean/proto
	rm -f go.sum

test: # -count=1 disables cache
	go test -v -race -count=1 .

.PHONY: all tidy clean test

include .make/lint.mk
include .make/proto.mk

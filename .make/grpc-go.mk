TOOLCHAIN ?= .tool
PROTOC_GEN_GO := $(TOOLCHAIN)/bin/protoc-gen-go

$(PROTOC_GEN_GO):
	GOPATH=$(shell pwd)/$(TOOLCHAIN) go install github.com/golang/protobuf/protoc-gen-go

%.pb.go: %.proto $(PROTOC) $(PROTOC_GEN_GO)
	PATH=$(shell pwd)/.tool/bin:$(PATH) protoc --go_out=plugins=grpc:. $<

clean/protoc-gen-go:
	rm -f $(PROTOC_GEN_GO)

.PHONY: clean/protoc-gen-go

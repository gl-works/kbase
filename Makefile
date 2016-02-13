.PHONY: compile doc fmt lint run test vendor_clean vendor_get vendor_update vet

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: fmt vet compile

compile:
	protoc --go_out=plugins=grpc:.  ./rpc/rpc.proto

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

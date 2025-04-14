.PHONY: fmt

all:

fmt:
	go fmt ./...

run-server:
	go run cmd/main.go server

generate-proto:
	rm -rf pkg/protocol
	mkdir pkg/protocol
	protoc.exe --go_out=pkg --go-grpc_out=pkg docs/protocol.proto
	protoc.exe --go_out=pkg --go-grpc_out=pkg docs/verifier.proto
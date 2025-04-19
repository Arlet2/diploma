.PHONY: fmt

all:

fmt:
	go fmt ./...

run-server:
	go run cmd/main.go server

run-ws:
	go run cmd/main.go ws

run-client:
	go run tools/client/client.go

generate-proto:
	rm -rf pkg/protocol
	mkdir pkg/protocol
	protoc.exe --go_out=pkg --go-grpc_out=pkg docs/protocol.proto
	protoc.exe --go_out=pkg --go-grpc_out=pkg docs/verifier.proto

docker-build:
	docker build -t push-diploma:1.0 .

auth-docker-build:
	docker build -f auth.Dockerfile -t auth-diploma:1.0 .
.PHONY: fmt

all:

fmt:
	go fmt ./...

run-server:
	go run cmd/main.go server
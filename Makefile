SHELL=/usr/bin/bash

debug:
	go run -v -trimpath -race ./cmd/api -config ./server.yaml

run:
	GIN_MODE=release go run -trimpath ./cmd/api -config ./server.yaml

test:
	go test -cover -covermode atomic -coverpkg=./... -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

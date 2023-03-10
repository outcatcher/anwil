SHELL=/usr/bin/bash

.PHONY: all build run vet test lint

build:
	go build -trimpath -o ./build/anwil ./cmd/api

vet:
	@go vet ./... && echo go vet OK

test: vet
	POSTGRES_PASSWORD= go test -cover -covermode atomic -coverpkg=./domains/... -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: vet
	go test -tags acceptance -covermode atomic -coverpkg=./domains/... -coverprofile=integration-coverage.out -timeout 1m -v ./...
	go tool cover -html=integration-coverage.out -o integration-coverage.html

lint:
	golangci-lint run ./...

local-deploy:
	POSTGRES_PORT=6634 docker compose up -d --build --wait


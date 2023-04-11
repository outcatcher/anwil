SHELL=/usr/bin/bash

.PHONY: all build run vet test lint

build:
	go build -trimpath -o ./build/anwil ./domains/api/cmd/server

vet:
	@go vet ./... && echo go vet OK

test: vet
	POSTGRES_PASSWORD= go test -cover -covermode atomic -coverpkg=./domains/... -coverprofile=coverage.out -count=1 -v ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: vet
	go test -tags acceptance -covermode atomic -coverpkg=./domains/... -coverprofile=coverage.out -timeout 1m -count=1 -v ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

local-deploy:
	POSTGRES_PORT=6634 docker compose up -d --build --wait


.DEFAULT_GOAL := build

fmt:
	goimports -l -w .
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o bin/main cmd/main.go
.PHONY:build

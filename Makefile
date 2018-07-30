.PHONY: build
SHELL := /bin/bash
export TESTS
header = "  \e[1;34m%-30s\e[m \n"
row = "\e[1mmake %-32s\e[m %-50s \n"

all:
	@printf $(header) "Build"
	@printf $(row) "build" "Build production binary."
	@printf $(row) "docker" "Build docker."
	@printf $(header) "Dev"
	@printf $(row) "run" "Run API in dev mode, all logging and race detector ON."
	@printf $(row) "test" "Run tests."
	@printf $(row) "vet" "Run go vet."
	@printf $(row) "lint" "Run gometalinter (you have to install it)."

build:
	go build

docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o refactored-octo-giggle
	docker build --no-cache -t refactored-octo-giggle .

run: 
	LOGXI=* go run -race main.go

test: 
	go test -count=1 -race -cover -v ./...

vet:
	go vet ./...

lint:
	gometalinter.v2 --disable=vetshadow --vendor ./...

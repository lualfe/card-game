.PHONY: help

help: ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test: ### run all tests
	go test ./...
.PHONY: test

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 ### run application
	go mod tidy && go mod download && go run cmd/app/main.go
.PHONY: run

build: ### builds the app image
	docker build --tag cards-deck .
.PHONY: build

run-build: ### runs the container app
	docker run --name cards-deck-app -p 8080:8080 cards-deck
.PHONY: run-build

build-and-run: build run-build ### build docker image and run
.PHONY: build-and-run
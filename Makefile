BINARY     := helloworld
IMAGE_NAME ?= helloworld
TAG        ?= latest

.PHONY: build test lint docker-build docker-push

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) ./cmd/helloworld

test:
	go test ./... -race -coverprofile=coverage.out

lint:
	go vet ./...
	staticcheck ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-push:
	docker push $(IMAGE_NAME):$(TAG)

.PHONY: build build-alpine

VERSION = $(shell cat VERSION)

IMAGE_NAME = git.furqan.io:5005/hjr265/bullet
IMAGE_TAG = $(VERSION)

build:
	go build -o bullet github.com/FurqanSoftware/bullet

build-alpine:
	docker run \
		-v `pwd`:/go/src/github.com/FurqanSoftware/bullet \
		-w /go/src/github.com/FurqanSoftware/bullet \
		golang:1.8-alpine \
		go build -o bullet github.com/FurqanSoftware/bullet

test:
	go test -v `go list ./... | grep -v /vendor/`

clean:
	go clean -i ./...

image: build-alpine
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: build build-alpine

IMAGE_NAME = registry.furqansoftware.net/tools/bullet
IMAGE_TAG = $(VERSION)

build:
	go build -o bullet github.com/FurqanSoftware/bullet

build-alpine:
	(cd foundry; docker build -t bullet-foundry .)
	docker run \
		-v `pwd`:/go/src/github.com/FurqanSoftware/bullet \
		-w /go/src/github.com/FurqanSoftware/bullet \
		bullet-foundry \
		go build -buildvcs=false -o bullet github.com/FurqanSoftware/bullet

test:
	go test -v `go list ./... | grep -v /vendor/`

clean:
	go clean -i ./...

image: build-alpine
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

IMAGE_NAME = registry.furqansoftware.net/tools/bullet
IMAGE_TAG = $(VERSION)

.PHONY: build
build:
	go build -o bullet github.com/FurqanSoftware/bullet

.PHONY: build.alpine
build.alpine:
	(cd foundry; docker build -t bullet-foundry .)
	docker run \
		-v `pwd`:/go/src/github.com/FurqanSoftware/bullet \
		-w /go/src/github.com/FurqanSoftware/bullet \
		bullet-foundry \
		go build -buildvcs=false -o bullet github.com/FurqanSoftware/bullet

.PHONY: build.dockerimage
build.dockerimage: build.alpine
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: lint.tools.install
lint.tools.install:
	go install honnef.co/go/tools/cmd/staticcheck@2023.1.2

.PHONY: clean
clean:
	go clean -i ./...


all: test build image

.PHONY: build
build:
	bin/build

build-image:
	bin/build-image

.PHONY: image
image: build build-image

vet:
	bin/vet

lint:
	bin/lint

test-unit:
	bin/test-unit

test-integration:
	bin/test-integration

test-e2e:
	bin/test-e2e

test: vet lint test-unit

tools:
	bin/tools

check-scripts:
	bin/check-scripts

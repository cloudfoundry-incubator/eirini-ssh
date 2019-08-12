all: test-unit build image

.PHONY: build
build:
	bin/build-ssh-proxy

export NAMESPACE ?= default

vet:
	bin/vet

lint:
	bin/lint

test-unit:
	bin/test-unit

test: vet lint test-unit

tools:
	bin/tools

check-scripts:
	bin/check-scripts

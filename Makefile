export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_PROXY_VERSION := latest

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-proxy ./cmd/onos-proxy

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps license_check linters
	go test -race github.com/onosproject/onos-proxy/...

jenkins-test: # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build deps license_check linters
	TEST_PACKAGES=github.com/onosproject/onos-proxy/pkg/... ./build/build-tools/build/jenkins/make-unit

onos-proxy-docker: # @HELP build onos-proxy base Docker image
	@go mod vendor
	docker build . -f build/onos-proxy/Dockerfile \
		-t onosproject/onos-proxy:${ONOS_PROXY_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-proxy-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-proxy:${ONOS_PROXY_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION} onosproject/onos-proxy

jenkins-publish: jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	./build/build-tools/release-merge-commit
	./build/build-tools/build/docs/push-docs

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-proxy/onos-proxy ./cmd/dummy/dummy


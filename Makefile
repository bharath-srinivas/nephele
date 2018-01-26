BINARY = aws_go
GOARCH = amd64
BUILD_DIR=$(shell pwd)

LDFLAGS = -ldflags "-s -w"

all: clean test build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build linux binaries.
linux:
	@cd cmd; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_linux_${GOARCH} .
.PHONY: linux

# Build release binaries.
build: linux
.PHONY: build

# Clean build artifacts.
clean:
	@rm -f ${BINARY}_linux_*
.PHONY: clean
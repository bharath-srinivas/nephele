BINARY = aws_go
BUILD_DIR=$(shell pwd)

LDFLAGS = -ldflags="-s -w"

all: clean test build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build linux binaries.
linux:
	@cd cmd/aws-go; \
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_linux_amd64 .
.PHONY: linux

linux-386:
	@cd cmd/aws-go; \
	GOOS=linux GOARCH=386 CGO_ENABLED=1 CFLAGS=-m32 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_linux_386 .
.PHONE: linux-386

# Build release binaries.
build: linux linux-386
.PHONY: build

# Clean build artifacts.
clean:
	@rm -f ${BINARY}_linux_*
.PHONY: clean
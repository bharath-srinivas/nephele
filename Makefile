BINARY = nephele
BUILD_DIR = $(shell pwd)
GO_VERSION = $(shell go version | cut -d ' ' -f 3 | sed -e 's/go//')
LDFLAGS = -ldflags="-s -w"

all: clean test build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build linux binaries.
linux:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=linux/amd64 ./cmd/nephele; \
	mv ${BINARY}-linux-amd64 ${BINARY}_linux_amd64
.PHONY: linux

linux-386:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=linux/386 ./cmd/nephele; \
	mv ${BINARY}-linux-386 ${BINARY}_linux_386
.PHONY: linux-386

# Build mac binaries.
darwin:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=darwin-10.9/amd64 ./cmd/nephele; \
	mv ${BINARY}-darwin-10.9-amd64 ${BINARY}_darwin_amd64
.PHONY: darwin

darwin-386:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=darwin-10.9/386 ./cmd/nephele; \
	mv ${BINARY}-darwin-10.9-386 ${BINARY}_darwin_386
.PHONY: darwin-386

# Build windows binaries.
windows:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=windows-6.1/amd64 ./cmd/nephele; \
	mv ${BINARY}-windows-6.1-amd64.exe ${BINARY}_windows_amd64.exe
.PHONY: windows

windows-386:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=windows-6.1/386 ./cmd/nephele; \
    mv ${BINARY}-windows-6.1-386.exe ${BINARY}_windows_386.exe
.PHONY: windows-386

# Build release binaries.
build:
	xgo -go ${GO_VERSION} ${LDFLAGS} -out ${BINARY} --targets=linux/amd64,linux/386,darwin-10.9/amd64,darwin-10.9/386,windows-6.1/amd64,windows-6.1/386 ./cmd/nephele; \
	mv ${BINARY}-linux-amd64 ${BINARY}_linux_amd64; \
	mv ${BINARY}-linux-386 ${BINARY}_linux_386; \
	mv ${BINARY}-darwin-10.9-amd64 ${BINARY}_darwin_amd64; \
	mv ${BINARY}-darwin-10.9-386 ${BINARY}_darwin_386; \
	mv ${BINARY}-windows-6.1-amd64.exe ${BINARY}_windows_amd64.exe; \
    mv ${BINARY}-windows-6.1-386.exe ${BINARY}_windows_386.exe
.PHONY: build

# Clean build artifacts.
clean:
	@rm -f ${BINARY}_linux_*
	@rm -f ${BINARY}_win*
	@rm -f ${BINARY}_darwin*
.PHONY: clean
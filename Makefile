BINARY = nephele
BUILD_DIR=$(shell pwd)

LDFLAGS = -ldflags="-s -w"

all: clean test build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build linux binaries.
linux:
	@cd cmd/nephele; \
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_linux_amd64 .
.PHONY: linux

linux-386:
	@cd cmd/nephele; \
	GOOS=linux GOARCH=386 CGO_ENABLED=1 CFLAGS=-m32 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_linux_386 .
.PHONY: linux-386

# Build windows binaries.
windows:
	@cd cmd/nephele; \
	GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_win64.exe .
.PHONY: windows

windows-386:
	@cd cmd/nephele; \
	GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}_win32.exe .
.PHONY: windows-386

# Build release binaries.
build: linux linux-386 windows windows-386
.PHONY: build

# Clean build artifacts.
clean:
	@rm -f ${BINARY}_linux_*
	@rm -f ${BINARY}_win*
.PHONY: clean
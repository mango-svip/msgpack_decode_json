
GOARCH=amd64

CGO_ENABLED=0

GOBUILD=go build

OUTPUT=./bin
usage:
	@echo "make init"
	@echo "make build"
	@echo "make clean"

init:
	mkdir -p bin
	go mod tidy
clean:
	rm -rf  ./bin/*
build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=windows GOARCH=${GOARCH} ${GOBUILD} -o ${OUTPUT}/msg_unpack.exe .
build-osx:
	CGO_ENABLED=${CGO_ENABLED} GOOS=darwin GOARCH=${GOARCH} ${GOBUILD} -o ${OUTPUT}/msg_unpack_osx  .
build-linux:
	CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=${GOARCH} ${GOBUILD} -o ${OUTPUT}/msg_unpack_linux   .
build-all: clean build build-osx build-linux


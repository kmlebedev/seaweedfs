.PHONY : build install test

BINARY = weed

SOURCE_DIR = .

all: install

install:
	cd weed; go install

full_install:
	cd weed; go install -tags "elastic gocdk sqlite ydb tikv"

test:
	cd weed; go test -tags "elastic gocdk sqlite ydb tikv" -v ./...

build:
	docker build --build-arg --no-cache -t chrislusf/seaweedfs:local -f docker/Dockerfile.build .
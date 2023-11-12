VERSION=$(shell git describe --tags --always)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: docker
docker:
	docker build -f ../../deploy/docker/harmoni/Dockerfile --build-arg VERSION=$(VERSION) -t harmoni ../../
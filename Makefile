BINARY     := parameters-iam
CMD        := ./cmd/server
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -X github.com/jasonmiller-cc/parameters-core/pkg/version.Version=$(VERSION) \
              -X github.com/jasonmiller-cc/parameters-core/pkg/version.Commit=$(COMMIT) \
              -X github.com/jasonmiller-cc/parameters-core/pkg/version.BuildTime=$(BUILD_TIME)

.PHONY: build run test lint vet tidy clean docker-build docker-run

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(CMD)

run: build
	./bin/$(BINARY)

test:
	go test ./... -race -coverprofile=coverage.out

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/ coverage.out

docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) \
	  -t parameters-iam:$(VERSION) -t parameters-iam:latest .

docker-run:
	docker run --rm -p 8080:8080 parameters-iam:latest

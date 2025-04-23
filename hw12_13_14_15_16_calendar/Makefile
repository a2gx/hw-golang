BIN := ./bin/calendar
DOCKER_IMG="calendar:develop"
CONFIG_FILE := "./configs/config.yaml"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN)_sender -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config $(CONFIG_FILE)

start:
	go run ./cmd/calendar -config $(CONFIG_FILE)

start-grpc:
	go run ./cmd/grpc_client -config $(CONFIG_FILE)

start-migration:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=up

start-migration-down:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=down

start-migration-new:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=create -name=new

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./cmd/...
	golangci-lint run ./internal/...
	golangci-lint run ./pkg/...

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative api/calendar.proto

dc-up:
	docker-compose up --build

.PHONY: build run start start-grpc-client start-migration-up start-migration-down start-migration-create start-initialize build-img run-img version test lint proto-gen dc-up

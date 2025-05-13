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

version: build
	$(BIN) version

start:
	go run ./cmd/calendar -config $(CONFIG_FILE)

start-scheduler:
	go run ./cmd/calendar_scheduler -config $(CONFIG_FILE)

start-sender:
	go run ./cmd/calendar_sender -config $(CONFIG_FILE)

start-grpc:
	go run ./cmd/grpc_client -config $(CONFIG_FILE)

start-migration:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=up

start-migration-down:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=down

start-migration-new:
	go run ./cmd/migration -config $(CONFIG_FILE) -command=create -name=new

test:
	go test -race ./internal/... ./pkg/...

integration-tests:
	docker-compose -f docker-compose.yml -f docker-compose.tests.yml down -v
	docker-compose --env-file .env -f docker-compose.yml -f docker-compose.tests.yml up --build -d
	docker logs -f integration_tests
	exit $$(docker wait integration_tests)

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

up:
	docker-compose --env-file .env up --build -d
down:
	docker-compose down
logs:
	docker-compose logs -f

.PHONY: build run start version start-scheduler start-sender start-grpc start-migration start-migration-up start-migration-down start-migration-new test integration-tests lint generate up down logs

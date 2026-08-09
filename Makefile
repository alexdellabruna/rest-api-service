APP_NAME ?= task-red-hat
GO ?= go
PORT ?= 8080
LOCAL_PORT ?= $(PORT)
# allows both ipv4 and ipv6 connections
LISTENING_ADDRESS ?= [::]

.PHONY: all build run test clean fmt vet tidy

all: tidy test run
all-docker: tidy test build-docker run-docker

build:
	$(GO) build -o bin/$(APP_NAME) ./...

build-docker:
	docker build -t $(APP_NAME) .

run:
	$(GO) run ./cmd

run-docker:
	docker run -p $(LOCAL_PORT):$(PORT) -e LISTENING_ADDRESS=$(LISTENING_ADDRESS) -e HTTP_PORT=$(PORT) $(APP_NAME)

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin

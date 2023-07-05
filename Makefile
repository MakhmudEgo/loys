
BUILD_PATH = build
SERVER_NAME = loys
SERVER_NAME_LINUX = loys-linux

.PHONY: all build build-debug build-linux build-linux-debug lint

all: build build-debug build-linux build-linux-debug

build:
	go build -o $(BUILD_PATH)/$(SERVER_NAME) cmd/loys/main.go

build-debug:
	go build -race -o $(BUILD_PATH)/$(SERVER_NAME)-debug cmd/loys/main.go

build-linux:
	GOOS=linux go build -o $(BUILD_PATH)/$(SERVER_NAME_LINUX) cmd/loys/main.go

build-linux-debug:
	GOOS=linux go build -race -o $(BUILD_PATH)/$(SERVER_NAME_LINUX)-debug cmd/loys/main.go

lint:
	type golangci-lint >/dev/null 2>&1 || {go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2 }
	golangci-lint run -v


clean:
	rm -rf $(BUILD_PATH)/*
	rm -rf db/data/*

run:
	docker-compose up --build

re: clean run


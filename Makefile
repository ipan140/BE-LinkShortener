APP_NAME=redirector
BUILD_DIR=bin
CMD_PATH=./cmd/api

.PHONY: all build run test clean docker-build docker-run

all: build

build:
	@echo "Building single executable binary..."
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test -v ./...

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f *.db *.log

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

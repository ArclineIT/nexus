.PHONY: build run test lint clean

APP_NAME := nexus
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v -race -count=1 ./...

lint:
	go vet ./...

# run with live reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

clean:
	rm -rf $(BUILD_DIR)


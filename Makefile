# Makefile for user-authentication Go project

BINARY_NAME := user-auth-service
BUILD_DIR := build
CMD_DIR := ./cmd
SERVICE_FILE := auth_service.go

.PHONY: all build clean run

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(SERVICE_FILE)

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

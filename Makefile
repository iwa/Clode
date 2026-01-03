.PHONY: help build run

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  run      - Run the application"

build:
	@echo "Building..."
	go build -o ./bin/Clode ./cmd/Clode
	@echo "Binary built at bin/Clode"

run:
	go run ./cmd/Clode

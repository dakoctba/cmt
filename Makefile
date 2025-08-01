.PHONY: build clean test install uninstall help

# Binary name
BINARY_NAME=cmt

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Install globally (requires sudo on macOS/Linux)
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/

# Uninstall
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Run the application
run: build
	./$(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  install   - Install globally (requires sudo)"
	@echo "  uninstall - Uninstall from system"
	@echo "  run       - Build and run the application"
	@echo "  help      - Show this help message"

# Default target
.DEFAULT_GOAL := build

.PHONY: build run clean generate test

# Build the application
build: generate
	go build -o bin/app cmd/main.go

# Generate templ templates
generate:
	templ generate

# Run the application
run: build
	./bin/app

# Run the application in development mode (with auto-reload)
dev:
	templ generate --watch &
	go run cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf data/

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...
	templ fmt .

# Lint code
lint:
	golangci-lint run

# Create data directory
init-data:
	mkdir -p data

# Setup the project
setup: deps generate init-data

# Docker build
docker-build:
	docker build -t home-automation .

# Docker run
docker-run:
	docker run -p 8080:8080 -v $(PWD)/data:/app/data home-automation
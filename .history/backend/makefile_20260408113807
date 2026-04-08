.PHONY: dev run build test test-cover docker-up docker-down clean

dev:
    @test -f .env || cp .env.example .env
    docker compose up -d postgres
    @sleep 3
    go run ./cmd/server

run:
    go run ./cmd/server

build:
    CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server ./cmd/server

test:
    go test ./internal/... -v -count=1 -race

test-cover:
    go test ./internal/... -coverprofile=coverage.out -count=1
    go tool cover -func=coverage.out

docker-up:
    @test -f .env || cp .env.example .env
    docker compose up -d --build

docker-down:
    docker compose down

clean:
    rm -rf bin/ coverage.out
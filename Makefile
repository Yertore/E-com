.PHONY: up down run-catalog build-catalog test-catalog lint

# Инфраструктура
up:
	docker compose up -d

down:
	docker compose down

# Catalog Service
run-catalog:
	cd catalog-service && go run ./cmd/server

build-catalog:
	cd catalog-service && go build -o ../bin/catalog-server ./cmd/server

test-catalog:
	cd catalog-service && go test ./...
.PHONY: run-catalog run-order run-payment run-notification \
        build-catalog test-catalog infra-up infra-down

infra-up:
	docker compose up -d

infra-down:
	docker compose down

run-catalog:
	cd catalog-service && go run ./cmd/server

run-order:
	cd order-service && go run ./cmd/server

run-payment:
	cd payment-service && go run ./cmd/server

run-notification:
	cd notification-service && go run ./cmd/server

build-catalog:
	cd catalog-service && go build -o ../bin/catalog-server ./cmd/server

test-catalog:
	cd catalog-service && go test ./...

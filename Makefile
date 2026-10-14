include .env
export

.PHONY: env-up env-down env-port-forward env-port-close \
		redis-up redis-down redis-port-forward redis-port-close\
        migrate-create migrate-up migrate-down migrate-action \
        build run

build:
	@go build -o ./build/todoapp ./cmd/todoapp

run: build env-up
	./build/todoapp

env-up:
	@docker compose up -d --wait todoapp-postgres todoapp-redis port-forwarder redis-port-forwarder

env-down:
	@docker compose stop redis-port-forwarder port-forwarder todoapp-redis todoapp-postgres

env-port-forward:
	@docker compose up -d --wait todoapp-postgres port-forwarder

env-port-close:
	@docker compose stop port-forwarder

redis-up:
	@docker compose up -d --wait todoapp-redis redis-port-forwarder

redis-down:
	@docker compose stop redis-port-forwarder todoapp-redis

redis-port-forward:
	@docker compose up -d --wait todoapp-redis redis-port-forwarder

redis-port-close:
	@docker compose stop redis-port-forwarder

migrate-create:
	@docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: env-up
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@migrate \
    	-path ./migrations \
    	-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:5433/$(POSTGRES_DB)?sslmode=disable \
    	"$(action)"


VALIDATE_DIR := $(shell go list -m -f "{{.Dir}}" github.com/envoyproxy/protoc-gen-validate)

proto:
	@protoc -I api/proto \
	  -I $(VALIDATE_DIR) \
	  --go_out=pkg/pb --go_opt=paths=source_relative \
	  --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	  --validate_out="lang=go,paths=source_relative:pkg/pb" \
	  api/proto/v1/*.proto
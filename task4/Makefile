SHELL := /bin/bash

.PHONY: up down restart logs ps build test gen proto sqlc tidy

up:
	docker compose up --build -d

down:
	docker compose down -v

restart: down up

logs:
	docker compose logs -f --tail=200

ps:
	docker compose ps

build:
	docker compose build

tidy:
	cd repo-stat && go mod tidy

gen: proto sqlc

proto:
	cd repo-stat && PATH="$$HOME/go/bin:$$PATH" protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/subscriber/subscriber.proto proto/collector/collector.proto proto/processor/processor.proto

sqlc:
	cd repo-stat && PATH="$$HOME/go/bin:$$PATH" sqlc generate -f subscriber/sqlc.yaml

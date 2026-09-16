include .env
export

.PHONY: up down logs db-shell migrate-up migrate-down migrate-create

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

db-shell:
	docker compose exec postgres psql -U school -d school_sms

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
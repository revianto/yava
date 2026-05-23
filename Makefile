.PHONY: dev dev-web dev-api db db-down migrate

dev:
	@make -j2 dev-api dev-web

dev-web:
	cd apps/web && npm run dev

dev-api:
	cd apps/api && air

db:
	docker compose up -d

db-down:
	docker compose down

migrate:
	cd apps/api && go run ./cmd/migrate up

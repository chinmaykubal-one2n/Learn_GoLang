.PHONY: run migrate-up migrate-down migrate-status migrate-force migrate-create migrate-create-table-add-phone

include .env

run:
	@go run main.go

migrate-up:
	@migrate -path internal/db/migrations -database $(DB_URL) up

migrate-down:
	@migrate -path internal/db/migrations -database $(DB_URL) down 1

migrate-status:
	@migrate -path internal/db/migrations -database $(DB_URL) version

migrate-force:
	@migrate -path internal/db/migrations -database $(DB_URL) force VERSION


	

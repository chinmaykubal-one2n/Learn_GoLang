.PHONY: run test migrate-diff migrate-apply migrate-down migrate-status swag

include .env

run:
	@go run main.go

test:
	@go test ./internal/tests/

migrate-diff:
	@atlas migrate diff --env gorm  

migrate-apply:
	@atlas migrate apply --env gorm --url $(DB_URL)

migrate-status:
	@atlas migrate status --env gorm --url $(DB_URL)   	

swag:
	@swag init -d ./,./internal/handler/,./internal/model/
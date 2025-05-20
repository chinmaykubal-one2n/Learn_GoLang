.PHONY: run unit-test integration-test migrate-diff migrate-apply migrate-down migrate-status swag

include .env

run:
	@go run main.go

unit-test:
	@go test ./internal/tests/unit-test/unit_test.go

integration-test:
	@go test ./internal/tests/integration-test/integration_test.go

migrate-diff:
	@atlas migrate diff --env gorm  

migrate-apply:
	@atlas migrate apply --env gorm --url $(DB_URL)

migrate-status:
	@atlas migrate status --env gorm --url $(DB_URL)   	

swag:
	@swag init -d ./,./internal/handler/,./internal/model/
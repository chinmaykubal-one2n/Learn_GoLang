.PHONY: run migrate-diff migrate-apply migrate-down migrate-status

include .env

run:
	@go run main.go

migrate-diff:
	@atlas migrate diff --env gorm  

migrate-apply:
	@atlas migrate apply --env gorm --url $(DB_URL)

migrate-status:
	@atlas migrate status --env gorm --url $(DB_URL)   	

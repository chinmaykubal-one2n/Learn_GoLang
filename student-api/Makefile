.PHONY: run migrate-diff migrate-apply migrate-down migrate-status

run:
	@go run main.go

migrate-diff:
	@atlas migrate diff --env gorm

migrate-apply:
	@atlas migrate apply --env gorm

migrate-down:
	@atlas migrate down --env gorm 

migrate-status:
	@atlas migrate status --env gorm 

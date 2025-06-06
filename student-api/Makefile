.PHONY: run handler-unit-test service-unit-test coverage integration-test migrate-diff migrate-apply migrate-down migrate-status swag seed-data docker-build docker-run docker-remove

include .env

run:
	@go run main.go

handler-unit-test:
	@go test ./internal/handler/student_handler_test.go

service-unit-test:
	@go test ./internal/service/student_service_test.go

coverage:
	@go test ./internal/handler ./internal/service -coverprofile=coverage.out
	@go tool cover -func=coverage.out

integration-test:
	@go test ./internal/integration/integration_test.go

migrate-diff:
	@atlas migrate diff --env gorm  

migrate-apply:
	@atlas migrate apply --env gorm --url $(DB_URL)

migrate-status:
	@atlas migrate status --env gorm --url $(DB_URL)   	

swag:
	@swag init -d ./,./internal/handler/,./internal/model/

seed-data:
	@go run ./internal/cmd/seeder/main.go

docker-build:
	@docker build -t student-api-go:multistage .

docker-run:
	@docker run -d -p 8090:8090 \
		--env-file .env.docker \
		--add-host=host.docker.internal:host-gateway \
		--network signoz-net \
		--name student-api-go \
		student-api-go:multistage

docker-remove:
	@docker rm -f student-api-go 
	@docker rmi -f student-api-go:multistage
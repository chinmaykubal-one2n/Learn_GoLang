.PHONY: init build test race coverage 

init:
	@test -f go.mod || go mod init wc-tool

build:
	@go build -o wc main.go

test:
	@go test 

race:
	@go run -race main.go ./shakespeare-db/*

coverage:
	@go test -coverprofile=coverage.out
	@go tool cover -func=coverage.out
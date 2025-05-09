.PHONY: init build test coverage race

init:
	@test -f go.mod ||  go mod init grep-cli  

build:
	@go build -o mygrep 

test:
	@go test ./cmd/

coverage:
	@go test -coverprofile=coverage.out ./cmd/
	@go tool cover -func=coverage.out  

race:
	@go run -race main.go -i -r hello ./test-data/
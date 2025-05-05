.PHONY: init build test

init:
	@test -f go.mod ||  go mod init grep-cli ./cmd/

build:
	@go build -o mygrep ./cmd/

test:
	@go test ./cmd/

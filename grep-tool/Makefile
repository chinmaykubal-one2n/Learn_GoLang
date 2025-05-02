.PHONY: init build test

init:
	@test -f go.mod || go mod init grep-cli

build:
	@go build -o mygrep 

test:
	@go test 

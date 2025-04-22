.PHONY: init build test

init:
	@test -f go.mod || go mod init wc-tool

build:
	@go build -o wc main.go

test:
	@go test 



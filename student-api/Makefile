.PHONY: run post get

run:
	@go run main.go

post:
	@curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{"name":"John","age":25,"email":"john@example.com"}'

get:
	@curl http://localhost:8080/students

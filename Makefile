build:
	@go build -o bin/image_resizer_b cmd/main.go

run: build
	@./bin/image_resizer_b cmd/main.go

test:
	@go test -v ./tests/handler_test.go

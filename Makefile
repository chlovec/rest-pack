build:
	@go build -o bin/example example/main.go

run: build
	@./bin/example

test:
	@go test -coverprofile=coverage.out -v ./...

test-html: test
    @go tool cover -html=coverage.out
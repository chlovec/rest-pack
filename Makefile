example-mysql:
	@go run examples/mysql/main.go

example-api:
	@go run examples/api/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

test-html: test
    @go tool cover -html=coverage.out
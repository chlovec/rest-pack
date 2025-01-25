example-mysql:
	@go run examples/mysql/main.go

example-api:
	@go run example_rest_api/cmd/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

test-html: test
    @go tool cover -html=coverage.out
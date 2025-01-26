run-example:
	@go run examples/cmd/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

test-html: test
    @go tool cover -html=coverage.out
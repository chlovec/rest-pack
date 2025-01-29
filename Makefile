run-example:
	@go run examples/cmd/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

test-rpt: test
	@go tool cover -html=coverage.out -o coverage.html
	@xdg-open coverage.html || open coverage.html
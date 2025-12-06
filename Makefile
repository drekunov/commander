lint:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	golangci-lint run

modup:
	go get -u ./...
	go mod tidy

lint:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	golangci-lint run

modup:
	go get -u ./...
	go mod tidy

run:
	go run ./cmd/commander/main.go

build:
	go build -o ./commander ./cmd/commander/main.go
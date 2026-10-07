run:
	go run ./cmd

build:
	go build -o ./bin/server ./cmd

test:
	go test ./internal/...
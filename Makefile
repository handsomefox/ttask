.DEFAULT_GOAL := run

test:
	CGO_ENABLED=1 go test -v -race ./...
.PHONY:test

run:
	go run ./cmd/api-server/main.go
.PHONY:run

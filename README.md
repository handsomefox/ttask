# ttask

## How to run

Either do:

```bash
make run
```

or

```bash
go run ./cmd/api-server/main.go
```

## How to run tests

Run:

```bash
make test
```

or

```bash
CGO_ENABLED=1 go test -v -race ./...
```

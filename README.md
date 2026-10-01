# Error Handling Lab

A small Go app demonstrating layered user lookup and wrapped errors. It uses only the Go standard library.

## Run

Run these commands from the project root, where `config.json` is located:

```sh
go run ./cmd/app 42
go run ./cmd/app 1
go run ./cmd/app abc
go run ./cmd/app 0
```

ID `42` demonstrates the not-found query and sentinel messages. ID `1` succeeds silently. Invalid IDs print an error to stderr and exit with status 1. Successful and not-found lookups write `audit.log` in the current directory.

## Verify

```sh
go build ./...
go vet ./...
go test ./...
```

GitHub Actions runs formatting, vet, tests, and build for pushes to `main` and pull requests.